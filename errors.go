// Copyright © 2015, 2016 Hraban Luyat <hraban@0brg.net>
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

type opusFileError int

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
		return "libopus error: %d (unknown code)"
	}
}

// opuserr translates libopus (not libopusfile) error codes to human readable
// strings
func opuserr(code int) error {
	return fmt.Errorf("opus: %s", C.GoString(C.opus_strerror(C.int(code))))
}
