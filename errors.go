// Copyright © 2015-2017 Go Opus Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

package opus

import (
	"fmt"
)

/*
#cgo pkg-config: opus opusfile
#ifdef __APPLE__
#include <opus.h>
#else
#include <opus/opus.h>
#endif
#include <opusfile.h>

// Access the preprocessor from CGO

// Errors for libopus
const int CONST_OPUS_OK = OPUS_OK;
const int CONST_OPUS_BAD_ARG = OPUS_BAD_ARG;
const int CONST_OPUS_BUFFER_TOO_SMALL = OPUS_BUFFER_TOO_SMALL;
const int CONST_OPUS_INTERNAL_ERROR = OPUS_INTERNAL_ERROR;
const int CONST_OPUS_INVALID_PACKET = OPUS_INVALID_PACKET;
const int CONST_OPUS_UNIMPLEMENTED = OPUS_UNIMPLEMENTED;
const int CONST_OPUS_INVALID_STATE = OPUS_INVALID_STATE;
const int CONST_OPUS_ALLOC_FAIL = OPUS_ALLOC_FAIL;

// Errors for libopusfile
const int CONST_OP_FALSE = OP_FALSE;
const int CONST_OP_EOF = OP_EOF;
const int CONST_OP_HOLE = OP_HOLE;
const int CONST_OP_EREAD = OP_EREAD;
const int CONST_OP_EFAULT = OP_EFAULT;
const int CONST_OP_EIMPL = OP_EIMPL;
const int CONST_OP_EINVAL = OP_EINVAL;
const int CONST_OP_ENOTFORMAT = OP_ENOTFORMAT;
const int CONST_OP_EBADHEADER = OP_EBADHEADER;
const int CONST_OP_EVERSION = OP_EVERSION;
const int CONST_OP_ENOTAUDIO = OP_ENOTAUDIO;
const int CONST_OP_EBADPACKET = OP_EBADPACKET;
const int CONST_OP_EBADLINK = OP_EBADLINK;
const int CONST_OP_ENOSEEK = OP_ENOSEEK;
const int CONST_OP_EBADTIMESTAMP = OP_EBADTIMESTAMP;
*/
import "C"

type Error int

var _ error = Error(0)

// Libopus errors
var (
	ErrOK             = Error(C.CONST_OPUS_OK)
	ErrBadArg         = Error(C.CONST_OPUS_BAD_ARG)
	ErrBufferTooSmall = Error(C.CONST_OPUS_BUFFER_TOO_SMALL)
	ErrInternalError  = Error(C.CONST_OPUS_INTERNAL_ERROR)
	ErrInvalidPacket  = Error(C.CONST_OPUS_INVALID_PACKET)
	ErrUnimplemented  = Error(C.CONST_OPUS_UNIMPLEMENTED)
	ErrInvalidState   = Error(C.CONST_OPUS_INVALID_STATE)
	ErrAllocFail      = Error(C.CONST_OPUS_ALLOC_FAIL)
)

// Error string (in human readable format) for libopus errors.
func (e Error) Error() string {
	return fmt.Sprintf("opus: %s", C.GoString(C.opus_strerror(C.int(e))))
}

type StreamError int

var _ error = StreamError(0)

// Libopusfile errors. The names are copied verbatim from the libopusfile
// library.
var (
	ErrStreamFalse        = StreamError(C.CONST_OP_FALSE)
	ErrStreamEOF          = StreamError(C.CONST_OP_EOF)
	ErrStreamHole         = StreamError(C.CONST_OP_HOLE)
	ErrStreamRead         = StreamError(C.CONST_OP_EREAD)
	ErrStreamFault        = StreamError(C.CONST_OP_EFAULT)
	ErrStreamImpl         = StreamError(C.CONST_OP_EIMPL)
	ErrStreamInval        = StreamError(C.CONST_OP_EINVAL)
	ErrStreamNotFormat    = StreamError(C.CONST_OP_ENOTFORMAT)
	ErrStreamBadHeader    = StreamError(C.CONST_OP_EBADHEADER)
	ErrStreamVersion      = StreamError(C.CONST_OP_EVERSION)
	ErrStreamNotAudio     = StreamError(C.CONST_OP_ENOTAUDIO)
	ErrStreamBadPacked    = StreamError(C.CONST_OP_EBADPACKET)
	ErrStreamBadLink      = StreamError(C.CONST_OP_EBADLINK)
	ErrStreamNoSeek       = StreamError(C.CONST_OP_ENOSEEK)
	ErrStreamBadTimestamp = StreamError(C.CONST_OP_EBADTIMESTAMP)
)

func (i StreamError) Error() string {
	switch i {
	case ErrStreamFalse:
		return "OP_FALSE"
	case ErrStreamEOF:
		return "OP_EOF"
	case ErrStreamHole:
		return "OP_HOLE"
	case ErrStreamRead:
		return "OP_EREAD"
	case ErrStreamFault:
		return "OP_EFAULT"
	case ErrStreamImpl:
		return "OP_EIMPL"
	case ErrStreamInval:
		return "OP_EINVAL"
	case ErrStreamNotFormat:
		return "OP_ENOTFORMAT"
	case ErrStreamBadHeader:
		return "OP_EBADHEADER"
	case ErrStreamVersion:
		return "OP_EVERSION"
	case ErrStreamNotAudio:
		return "OP_ENOTAUDIO"
	case ErrStreamBadPacked:
		return "OP_EBADPACKET"
	case ErrStreamBadLink:
		return "OP_EBADLINK"
	case ErrStreamNoSeek:
		return "OP_ENOSEEK"
	case ErrStreamBadTimestamp:
		return "OP_EBADTIMESTAMP"
	default:
		return "libopusfile error: %d (unknown code)"
	}
}
