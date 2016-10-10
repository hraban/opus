// Copyright © 2015, 2016 Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

package opus

import (
	"fmt"
)

/*
#cgo pkg-config: opus opusfile
#include <opus/opus.h>
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

type opusError int

var _ error = opusError(0)

// Libopus errors
var (
	ERR_OPUS_OK               = opusError(C.CONST_OPUS_OK)
	ERR_OPUS_BAD_ARG          = opusError(C.CONST_OPUS_BAD_ARG)
	ERR_OPUS_BUFFER_TOO_SMALL = opusError(C.CONST_OPUS_BUFFER_TOO_SMALL)
	ERR_OPUS_INTERNAL_ERROR   = opusError(C.CONST_OPUS_INTERNAL_ERROR)
	ERR_OPUS_INVALID_PACKET   = opusError(C.CONST_OPUS_INVALID_PACKET)
	ERR_OPUS_UNIMPLEMENTED    = opusError(C.CONST_OPUS_UNIMPLEMENTED)
	ERR_OPUS_INVALID_STATE    = opusError(C.CONST_OPUS_INVALID_STATE)
	ERR_OPUS_ALLOC_FAIL       = opusError(C.CONST_OPUS_ALLOC_FAIL)
)

// Error string (in human readable format) for libopus errors.
func (e opusError) Error() string {
	return fmt.Sprintf("opus: %s", C.GoString(C.opus_strerror(C.int(e))))
}

type opusFileError int

var _ error = opusFileError(0)

// Libopusfile errors. The names are copied verbatim from the libopusfile
// library.
var (
	ERR_OP_FALSE         = opusFileError(C.CONST_OP_FALSE)
	ERR_OP_EOF           = opusFileError(C.CONST_OP_EOF)
	ERR_OP_HOLE          = opusFileError(C.CONST_OP_HOLE)
	ERR_OP_EREAD         = opusFileError(C.CONST_OP_EREAD)
	ERR_OP_EFAULT        = opusFileError(C.CONST_OP_EFAULT)
	ERR_OP_EIMPL         = opusFileError(C.CONST_OP_EIMPL)
	ERR_OP_EINVAL        = opusFileError(C.CONST_OP_EINVAL)
	ERR_OP_ENOTFORMAT    = opusFileError(C.CONST_OP_ENOTFORMAT)
	ERR_OP_EBADHEADER    = opusFileError(C.CONST_OP_EBADHEADER)
	ERR_OP_EVERSION      = opusFileError(C.CONST_OP_EVERSION)
	ERR_OP_ENOTAUDIO     = opusFileError(C.CONST_OP_ENOTAUDIO)
	ERR_OP_EBADPACKET    = opusFileError(C.CONST_OP_EBADPACKET)
	ERR_OP_EBADLINK      = opusFileError(C.CONST_OP_EBADLINK)
	ERR_OP_ENOSEEK       = opusFileError(C.CONST_OP_ENOSEEK)
	ERR_OP_EBADTIMESTAMP = opusFileError(C.CONST_OP_EBADTIMESTAMP)
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
		return "libopusfile error: %d (unknown code)"
	}
}
