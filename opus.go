// Copyright © 2015, 2016 Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

package opus

/*
// Link opus using pkg-config.
#cgo pkg-config: opus
#include <opus/opus.h>

// Access the preprocessor from CGO
const int CONST_APPLICATION_VOIP = OPUS_APPLICATION_VOIP;
const int CONST_APPLICATION_AUDIO = OPUS_APPLICATION_AUDIO;
const int CONST_APPLICATION_RESTRICTED_LOWDELAY = OPUS_APPLICATION_RESTRICTED_LOWDELAY;
*/
import "C"

const (
	// Optimize encoding for VOIP
	APPLICATION_VOIP = C.OPUS_APPLICATION_VOIP
	// Optimize encoding for non-voice signals like music
	APPLICATION_AUDIO = C.OPUS_APPLICATION_AUDIO
	// Optimize encoding for low latency applications
	APPLICATION_RESTRICTED_LOWDELAY = C.OPUS_APPLICATION_RESTRICTED_LOWDELAY
	// Auto bitrate
	OPUS_AUTO = C.OPUS_AUTO
	// Maximum bitrate
	OPUS_BITRATE_MAX = C.OPUS_BITRATE_MAX
	//4 kHz passband
	OPUS_BANDWIDTH_NARROWBAND = C.OPUS_BANDWIDTH_NARROWBAND
	// 6 kHz passband
	OPUS_BANDWIDTH_MEDIUMBAND = C.OPUS_BANDWIDTH_MEDIUMBAND
	// 8 kHz passband
	OPUS_BANDWIDTH_WIDEBAND = C.OPUS_BANDWIDTH_WIDEBAND
	// 12 kHz passband
	OPUS_BANDWIDTH_SUPERWIDEBAND = C.OPUS_BANDWIDTH_SUPERWIDEBAND
	// 20 kHz passband
	OPUS_BANDWIDTH_FULLBAND = C.OPUS_BANDWIDTH_FULLBAND
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
