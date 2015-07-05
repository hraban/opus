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
	if err != nil || n == 0 {
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
		return fmt.Errorf("Reader function must be non-nil")
	}
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
	s.buf = make([]byte, maxEncodedFrameSize)
	return nil
}

func (s *Stream) Read() ([]int16, error) {
	pcm := make([]int16, xMAX_FRAME_SIZE)
	n := C.op_read(
		s.oggfile,
		(*C.opus_int16)(&pcm[0]),
		C.int(cap(pcm)),
		nil)
	if n < 0 {
		return nil, opusfileerr(n)
	}
	return pcm[:n], nil
}
