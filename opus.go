// Copyright © 2015, 2016 Hraban Luyat <hraban@0brg.net>
//
// License for use of this code is detailed in the LICENSE file

package opus

/*
// Statically link libopus. Requires a libopus.a in every directory you use this
// as a dependency in. Not great, but CGO doesn't offer anything better, right
// now. Unless you require everyone who USES this package to have libopus
// installed system-wide, which is more of a chore because it's so new. Everyone
// will end up having to build it from source anyway, might as well just dump
// the pre-built lib in here. At least it will be up to the package maintainer,
// not the user.
//
// If I missed something, and somebody knows a better way: please let me know.
#cgo pkg-config: opus
#include <opus/opus.h>

// Access the preprocessor from CGO
const int CONST_APPLICATION_VOIP = OPUS_APPLICATION_VOIP;
const int CONST_APPLICATION_AUDIO = OPUS_APPLICATION_AUDIO;
const int CONST_APPLICATION_RESTRICTED_LOWDELAY = OPUS_APPLICATION_RESTRICTED_LOWDELAY;
*/
import "C"

type Application int

// These variables should be constants, but for interoperability with CGO
// they're var. Don't change them, though!
var (
	// Optimize encoding for VOIP
	APPLICATION_VOIP = Application(C.CONST_APPLICATION_VOIP)
	// Optimize encoding for non-voice signals like music
	APPLICATION_AUDIO = Application(C.CONST_APPLICATION_AUDIO)
	// Optimize encoding for low latency applications
	APPLICATION_RESTRICTED_LOWDELAY = Application(C.CONST_APPLICATION_RESTRICTED_LOWDELAY)
)

const (
	xMAX_BITRATE       = 48000
	xMAX_FRAME_SIZE_MS = 60
	xMAX_FRAME_SIZE    = xMAX_BITRATE * xMAX_FRAME_SIZE_MS / 1000
	// Maximum size of an encoded frame. I actually have no idea, but this
	// looks like it's big enough.
	maxEncodedFrameSize = 10000
)

func Version() string {
	return C.GoString(C.opus_get_version_string())
}
