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
	if pcm == nil || len(pcm) == 0 {
		return nil, fmt.Errorf("opus: no data supplied")
	}
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
