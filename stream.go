package opus

import (
	"fmt"
	"io"
	"unsafe"
)

/*
#cgo CFLAGS: -std=c99 -Wall -Werror -pedantic -Ibuild/include/opus
#include <opusfile.h>
#include <string.h>
#include "callback_proxy.h"
*/
import "C"

type Stream struct {
	oggfile *C.OggOpusFile
	read    io.Reader
	// Preallocated buffer to pass to the reader
	buf []byte
}

//export go_readproxy
func go_readproxy(p unsafe.Pointer, cbuf *C.uchar, cmaxbytes C.int) C.int {
	stream := (*Stream)(p)
	maxbytes := int(cmaxbytes)
	if maxbytes > cap(stream.buf) {
		maxbytes = cap(stream.buf)
	}
	// Don't bother cleaning up old data because that's not required by the
	// io.Reader API.
	n, err := stream.read.Read(stream.buf[:maxbytes])
	if (err != nil && err != io.EOF) || n == 0 {
		return 0
	}
	C.memcpy(unsafe.Pointer(cbuf), unsafe.Pointer(&stream.buf[0]), C.size_t(n))
	return C.int(n)
}

var callbacks = C.struct_OpusFileCallbacks{
	read:  C.op_read_func(C.c_readproxy),
	seek:  nil,
	tell:  nil,
	close: nil,
}

func NewStream(read io.Reader) (*Stream, error) {
	var s Stream
	err := s.Init(read)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Init initializes a stream with an io.Reader to fetch opus encoded data from
// on demand. Errors from the reader are all transformed to an EOF, any actual
// error information is lost. The same happens when a read returns succesfully,
// but with zero bytes.
func (s *Stream) Init(read io.Reader) error {
	if s.oggfile != nil {
		return fmt.Errorf("opus stream is already initialized")
	}
	if read == nil {
		return fmt.Errorf("Reader must be non-nil")
	}
	s.read = read
	s.buf = make([]byte, maxEncodedFrameSize)
	var errno C.int
	oggfile := C.op_open_callbacks(
		unsafe.Pointer(s),
		&callbacks,
		nil,
		0,
		&errno)
	if errno != 0 {
		return opusfileerr(errno)
	}
	s.oggfile = oggfile
	return nil
}

// Read a chunk of raw opus data from the stream and decode it. Returns the
// number of decoded samples per channel. This means that a dual channel
// (stereo) feed will have twice as many samples as the value returned.
//
// Read may successfully read less bytes than requested, but it will never read
// exactly zero bytes succesfully if a non-zero buffer is supplied.
//
// The number of channels in the output data must be known in advance. It is
// possible to extract this information from the stream itself, but I'm not
// motivated to do that. Feel free to send a pull request.
func (s *Stream) Read(pcm []int16) (int, error) {
	if s.oggfile == nil {
		return 0, fmt.Errorf("opus stream is uninitialized or already closed")
	}
	if len(pcm) == 0 {
		return 0, nil
	}
	n := C.op_read(
		s.oggfile,
		(*C.opus_int16)(&pcm[0]),
		C.int(len(pcm)),
		nil)
	if n < 0 {
		return 0, opusfileerr(n)
	}
	if n == 0 {
		return 0, io.EOF
	}
	return int(n), nil
}

func (s *Stream) ReadFloat32(pcm []float32) (int, error) {
	if s.oggfile == nil {
		return 0, fmt.Errorf("opus stream is uninitialized or already closed")
	}
	if len(pcm) == 0 {
		return 0, nil
	}
	n := C.op_read_float(
		s.oggfile,
		(*C.float)(&pcm[0]),
		C.int(len(pcm)),
		nil)
	if n < 0 {
		return 0, opusfileerr(n)
	}
	if n == 0 {
		return 0, io.EOF
	}
	return int(n), nil
}

func (s *Stream) Close() error {
	if s.oggfile == nil {
		return fmt.Errorf("opus stream is uninitialized or already closed")
	}
	C.op_free(s.oggfile)
	if closer, ok := s.read.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
