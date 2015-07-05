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

var errDecUninitialized = fmt.Errorf("opus decoder uninitialized")

type Decoder struct {
	p *C.struct_OpusDecoder
	// Same purpose as encoder struct
	mem         []byte
	sample_rate int
}

// NewEncoder allocates a new Opus decoder and initializes it with the
// appropriate parameters. All related memory is managed by the Go GC.
func NewDecoder(sample_rate int, channels int) (*Decoder, error) {
	var dec Decoder
	err := dec.Init(sample_rate, channels)
	if err != nil {
		return nil, err
	}
	return &dec, nil
}

func (dec *Decoder) Init(sample_rate int, channels int) error {
	if dec.p != nil {
		return fmt.Errorf("opus decoder already initialized")
	}
	if channels != 1 && channels != 2 {
		return fmt.Errorf("Number of channels must be 1 or 2: %d", channels)
	}
	size := C.opus_decoder_get_size(C.int(channels))
	dec.sample_rate = sample_rate
	dec.mem = make([]byte, size)
	dec.p = (*C.OpusDecoder)(unsafe.Pointer(&dec.mem[0]))
	errno := C.opus_decoder_init(
		dec.p,
		C.opus_int32(sample_rate),
		C.int(channels))
	if errno != 0 {
		return opuserr(int(errno))
	}
	return nil
}

func (dec *Decoder) Decode(data []byte) ([]int16, error) {
	if dec.p == nil {
		return nil, errDecUninitialized
	}
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
	if dec.p == nil {
		return nil, errDecUninitialized
	}
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
