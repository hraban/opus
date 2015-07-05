// Copyright © 2015 Hraban Luyat <hraban@0brg.net>
//
// License for use of this code is detailed in the LICENSE file

package opus

import (
	"fmt"
	"unsafe"
)

/*
#cgo CFLAGS: -std=c99 -Wall -Werror -pedantic -Ibuild/include
#include <opus/opus.h>
*/
import "C"

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

func (dec *Decoder) Decode(data []byte) ([]int16, error) {
	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("opus: no data supplied")
	}
	// I don't know how big this frame will be, but this is the limit
	pcm := make([]int16, xMAX_FRAME_SIZE_MS*dec.sample_rate/1000)
	n := int(C.opus_decode(
		dec.p,
		(*C.uchar)(&data[0]),
		C.opus_int32(len(data)),
		(*C.opus_int16)(&pcm[0]),
		C.int(cap(pcm)),
		0))
	if n < 0 {
		return nil, opuserr(n)
	}
	return pcm[:n], nil
}

func (dec *Decoder) DecodeFloat32(data []byte) ([]float32, error) {
	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("opus: no data supplied")
	}
	pcm := make([]float32, xMAX_FRAME_SIZE_MS*dec.sample_rate/1000)
	n := int(C.opus_decode_float(
		dec.p,
		(*C.uchar)(&data[0]),
		C.opus_int32(len(data)),
		(*C.float)(&pcm[0]),
		C.int(cap(pcm)),
		0))
	if n < 0 {
		return nil, opuserr(n)
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
