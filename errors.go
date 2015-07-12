package opus

import (
	"fmt"
)

/*
#cgo CFLAGS: -std=c99 -Wall -Werror -pedantic -Ibuild/include
#include <opus/opus.h>
*/
import "C"

type opusFileError int

var opusfileErrcodes = map[C.int]error{}

const (
	ERR_OP_FALSE         opusFileError = -1
	ERR_OP_EOF                         = -2
	ERR_OP_HOLE                        = -3
	ERR_OP_EREAD                       = -128
	ERR_OP_EFAULT                      = -129
	ERR_OP_EIMPL                       = -130
	ERR_OP_EINVAL                      = -131
	ERR_OP_ENOTFORMAT                  = -132
	ERR_OP_EBADHEADER                  = -133
	ERR_OP_EVERSION                    = -134
	ERR_OP_ENOTAUDIO                   = -135
	ERR_OP_EBADPACKET                  = -136
	ERR_OP_EBADLINK                    = -137
	ERR_OP_ENOSEEK                     = -138
	ERR_OP_EBADTIMESTAMP               = -139
)

func (i opusFileError) Error() string {
	switch i {
	case ERR_OP_FALSE:
		return "OP_FALSE"
	case ERR_OP_EOF:
		return "OP_EOF"
	case ERR_OP_HOLE:
		return "OP_HOLE"
	case ERR_OP_EREAD:
		return "OP_EREAD"
	case ERR_OP_EFAULT:
		return "OP_EFAULT"
	case ERR_OP_EIMPL:
		return "OP_EIMPL"
	case ERR_OP_EINVAL:
		return "OP_EINVAL"
	case ERR_OP_ENOTFORMAT:
		return "OP_ENOTFORMAT"
	case ERR_OP_EBADHEADER:
		return "OP_EBADHEADER"
	case ERR_OP_EVERSION:
		return "OP_EVERSION"
	case ERR_OP_ENOTAUDIO:
		return "OP_ENOTAUDIO"
	case ERR_OP_EBADPACKET:
		return "OP_EBADPACKET"
	case ERR_OP_EBADLINK:
		return "OP_EBADLINK"
	case ERR_OP_ENOSEEK:
		return "OP_ENOSEEK"
	case ERR_OP_EBADTIMESTAMP:
		return "OP_EBADTIMESTAMP"
	default:
		return "libopus error: %d (unknown code)"
	}
}

// opuserr translates libopus (not libopusfile) error codes to human readable
// strings
func opuserr(code int) error {
	return fmt.Errorf("opus: %s", C.GoString(C.opus_strerror(C.int(code))))
}
