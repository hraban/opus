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

type Encoder struct {
	p *C.struct_OpusEncoder
	// Memory for the encoder struct allocated on the Go heap to allow Go GC to
	// manage it (and obviate need to free())
	mem []byte
}

// NewEncoder allocates a new Opus encoder and initializes it with the
// appropriate parameters. All related memory is managed by the Go GC.
func NewEncoder(sample_rate int, channels int, application Application) (*Encoder, error) {
	if channels != 1 && channels != 2 {
		return nil, fmt.Errorf("Number of channels must be 1 or 2: %d", channels)
	}
	var e Encoder
	size := C.opus_encoder_get_size(C.int(channels))
	e.mem = make([]byte, size)
	e.p = (*C.OpusEncoder)(unsafe.Pointer(&e.mem[0]))
	errno := int(C.opus_encoder_init(
		e.p,
		C.opus_int32(sample_rate),
		C.int(channels),
		C.int(application)))
	if errno != 0 {
		return nil, opuserr(int(errno))
	}
	return &e, nil
}

func (enc *Encoder) Encode(pcm []int16) ([]byte, error) {
	if pcm == nil || len(pcm) == 0 {
		return nil, fmt.Errorf("opus: no data supplied")
	}
	// I never know how much to allocate
	data := make([]byte, 10000)
	n := int(C.opus_encode(
		enc.p,
		(*C.opus_int16)(&pcm[0]),
		C.int(len(pcm)),
		(*C.uchar)(&data[0]),
		C.opus_int32(cap(data))))
	if n < 0 {
		return nil, opuserr(n)
	}
	return data[:n], nil
}

func (enc *Encoder) EncodeFloat32(pcm []float32) ([]byte, error) {
	if pcm == nil || len(pcm) == 0 {
		return nil, fmt.Errorf("opus: no data supplied")
	}
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
