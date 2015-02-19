package opus

import (
	"fmt"
	"unsafe"
)

/*
#cgo CFLAGS: -std=c99 -Wall -Werror -pedantic
// Statically link libopus
#cgo LDFLAGS: libopusbuild/lib/libopus.a -lm
#include <opus/opus.h>
*/
import "C"

type Application int

// TODO: Get from lib because #defines can change
const APPLICATION_VOIP Application = 2048
const APPLICATION_AUDIO Application = 2049
const APPLICATION_RESTRICTED_LOWDELAY Application = 2051

const xMAX_BITRATE = 48000
const xMAX_FRAME_SIZE_MS = 60
const xMAX_FRAME_SIZE = xMAX_BITRATE * xMAX_FRAME_SIZE_MS / 1000

func Version() string {
	return C.GoString(C.opus_get_version_string())
}

type Encoder struct {
	p *C.struct_OpusEncoder
}

func opuserr(code int) error {
	return fmt.Errorf("opus: %s", C.GoString(C.opus_strerror(C.int(code))))
}

func NewEncoder(sample_rate int, channels int, application Application) (*Encoder, error) {
	var errno int
	p := C.opus_encoder_create(C.opus_int32(sample_rate), C.int(channels), C.int(application), (*C.int)(unsafe.Pointer(&errno)))
	if errno != 0 {
		return nil, opuserr(errno)
	}
	return &Encoder{p: p}, nil
}

func (enc *Encoder) EncodeFloat32(pcm []float32) ([]byte, error) {
	// I never know how much to allocate
	data := make([]byte, 10000)
	n := int(C.opus_encode_float(
		enc.p,
		(*C.float)(&pcm[0]),
		C.int(len(pcm)),
		(*C.uchar)(&data[0]),
		C.opus_int32(cap(data))))
	if n < 0 {
		return nil, opuserr(n)
	}
	return data[:n], nil
}

// Returns an error if the encoder was already closed
func (enc *Encoder) Close() error {
	if enc.p == nil {
		return fmt.Errorf("opus: encoder already closed")
	}
	C.opus_encoder_destroy(enc.p)
	enc.p = nil
	return nil
}

type Decoder struct {
	p           *C.struct_OpusDecoder
	sample_rate int
}

func NewDecoder(sample_rate int, channels int) (*Decoder, error) {
	var errno int
	p := C.opus_decoder_create(C.opus_int32(sample_rate), C.int(channels), (*C.int)(unsafe.Pointer(&errno)))
	if errno != 0 {
		return nil, opuserr(errno)
	}
	dec := &Decoder{
		p:           p,
		sample_rate: sample_rate,
	}
	return dec, nil
}

func (dec *Decoder) DecodeFloat32(data []byte) ([]float32, error) {
	var errno int
	// I don't know how big this frame will be, but this is the limit
	pcm := make([]float32, xMAX_FRAME_SIZE_MS*dec.sample_rate/1000)
	n := int(C.opus_decode_float(
		dec.p,
		(*C.uchar)(&data[0]),
		C.opus_int32(len(data)),
		(*C.float)(&pcm[0]),
		C.int(cap(pcm)),
		0))
	if n < 0 {
		return nil, opuserr(errno)
	}
	return pcm[:n], nil
}

// Returns an error if the encoder was already closed
func (dec *Decoder) Close() error {
	if dec.p == nil {
		return fmt.Errorf("opus: decoder already closed")
	}
	C.opus_decoder_destroy(dec.p)
	dec.p = nil
	return nil
}
