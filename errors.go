package opus

import (
	"errors"
	"fmt"
)

/*
#cgo CFLAGS: -std=c99 -Wall -Werror -pedantic -Ibuild/include
#include <opus/opus.h>
*/
import "C"

var opusfileErrcodes = map[C.int]error{
	-1:   errors.New("OP_FALSE"),
	-2:   errors.New("OP_EOF"),
	-3:   errors.New("OP_HOLE"),
	-128: errors.New("OP_EREAD"),
	-129: errors.New("OP_EFAULT"),
	-130: errors.New("OP_EIMPL"),
	-131: errors.New("OP_EINVAL"),
	-132: errors.New("OP_ENOTFORMAT"),
	-133: errors.New("OP_EBADHEADER"),
	-134: errors.New("OP_EVERSION"),
	-135: errors.New("OP_ENOTAUDIO"),
	-136: errors.New("OP_EBADPACKET"),
	-137: errors.New("OP_EBADLINK"),
	-138: errors.New("OP_ENOSEEK"),
	-139: errors.New("OP_EBADTIMESTAMP"),
}

// opusfileerr maps libopusfile error codes to human readable strings
func opusfileerr(code C.int) error {
	err, ok := opusfileErrcodes[code]
	if ok {
		return err
	}
	return fmt.Errorf("libopus error: %d (unknown code)", int(code))
}

// opuserr translates libopus (not libopusfile) error codes to human readable
// strings
func opuserr(code int) error {
	return fmt.Errorf("opus: %s", C.GoString(C.opus_strerror(C.int(code))))
}
