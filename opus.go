// Copyright © 2015-2017 Go Opus Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

package opus

/*
// Link opus using pkg-config.
#cgo pkg-config: opus
#ifdef __APPLE__
#include <opus.h>
#else
#include <opus/opus.h>
#endif

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
	// Optimize encoding for VoIP
	AppVoIP = Application(C.CONST_APPLICATION_VOIP)
	// Optimize encoding for non-voice signals like music
	AppAudio = Application(C.CONST_APPLICATION_AUDIO)
	// Optimize encoding for low latency applications
	AppRestrictedLowdelay = Application(C.CONST_APPLICATION_RESTRICTED_LOWDELAY)
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
