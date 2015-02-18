package opus

import (
	"fmt"
	"unsafe"
)

/*
#cgo CFLAGS: -std=c99 -Wall -Werror -pedantic
// Statically link libopus
#cgo LDFLAGS: libopusbuild/lib/libopus.a -lm
#include <stdio.h>
#include <stdlib.h>
#include <opus/opus.h>

int
engavg(const float *ar, int numfloats, char *buf, int bufsize)
{
    double sum = 0;

    for (int i = 0; i < numfloats; i++) {
        sum += ar[i];
    }

    int n = snprintf(buf, bufsize, "%g", sum / numfloats);
    if (n < 0 || n >= bufsize) {
        return -1;
    }

    return n;
}

*/
import "C"

func engavg(ar []float32) string {
	buf := make([]byte, 30)
	n := int(C.engavg(
		(*C.float)(&ar[0]),
		C.int(len(ar)),
		(*C.char)(unsafe.Pointer(&buf[0])),
		C.int(cap(buf))))
	if n < 0 {
		return fmt.Sprintf("Error: %d", n)
	}
	return string(buf[:n])

}

type Application int

// TODO: Get from lib because #defines can change
const APPLICATION_VOIP Application = 2048
const APPLICATION_AUDIO Application = 2049
const APPLICATION_RESTRICTED_LOWDELAY Application = 2051

type Encoder struct {
	p unsafe.Pointer
}

func NewEncoder(sample_rate int, channels int, application Application) (*Encoder, error) {
	var errno int
	p := unsafe.Pointer(C.opus_encoder_create(C.opus_int32(sample_rate), C.int(channels), C.int(application), (*C.int)(unsafe.Pointer(&errno))))
	if errno != 0 {
		return nil, fmt.Errorf("opus: %s", C.GoString(C.opus_strerror(C.int(errno))))
	}
	return &Encoder{p: p}, nil
}

func Version() string {
	return C.GoString(C.opus_get_version_string())
}
