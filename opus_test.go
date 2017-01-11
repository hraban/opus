// Copyright © 2015-2017 Go Opus Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

package opus

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	if ver := Version(); !strings.HasPrefix(ver, "libopus") {
		t.Errorf("Unexpected linked libopus version: " + ver)
	}
}

func TestOpusErrstr(t *testing.T) {
	// I scooped this -1 up from opus_defines.h, it's OPUS_BAD_ARG. Not pretty,
	// but it's better than not testing at all. Again, accessing #defines from
	// CGO is not possible.
	if ErrBadArg.Error() != "opus: invalid argument" {
		t.Errorf("Expected \"invalid argument\" error message for error code %d: %v",
			ErrBadArg, ErrBadArg)
	}
}

func TestCodec(t *testing.T) {
	// Create bogus input sound
	const G4 = 391.995
	const SAMPLE_RATE = 48000
	const FRAME_SIZE_MS = 60
	const FRAME_SIZE = SAMPLE_RATE * FRAME_SIZE_MS / 1000
	pcm := make([]int16, FRAME_SIZE)
	enc, err := NewEncoder(SAMPLE_RATE, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Fatalf("Error creating new encoder: %v", err)
	}
	addSine(pcm, SAMPLE_RATE, G4)
	data := make([]byte, 1000)
	n, err := enc.Encode(pcm, data)
	if err != nil {
		t.Fatalf("Couldn't encode data: %v", err)
	}
	data = data[:n]
	dec, err := NewDecoder(SAMPLE_RATE, 1)
	if err != nil || dec == nil {
		t.Fatalf("Error creating new decoder: %v", err)
	}
	n, err = dec.Decode(data, pcm)
	if err != nil {
		t.Fatalf("Couldn't decode data: %v", err)
	}
	if len(pcm) != n {
		t.Fatalf("Length mismatch: %d samples in, %d out", len(pcm), n)
	}
	// Checking the output programmatically is seriously not easy. I checked it
	// by hand and by ear, it looks fine. I'll just assume the API faithfully
	// passes error codes through, and that's that.
}

func TestCodecFloat32(t *testing.T) {
	// Create bogus input sound
	const G4 = 391.995
	const SAMPLE_RATE = 48000
	const FRAME_SIZE_MS = 60
	const FRAME_SIZE = SAMPLE_RATE * FRAME_SIZE_MS / 1000
	pcm := make([]float32, FRAME_SIZE)
	enc, err := NewEncoder(SAMPLE_RATE, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Fatalf("Error creating new encoder: %v", err)
	}
	addSineFloat32(pcm, SAMPLE_RATE, G4)
	data := make([]byte, 1000)
	n, err := enc.EncodeFloat32(pcm, data)
	if err != nil {
		t.Fatalf("Couldn't encode data: %v", err)
	}
	dec, err := NewDecoder(SAMPLE_RATE, 1)
	if err != nil || dec == nil {
		t.Fatalf("Error creating new decoder: %v", err)
	}
	// TODO: Uh-oh.. it looks like I forgot to put a data = data[:n] here, yet
	// the test is not failing. Why?
	n, err = dec.DecodeFloat32(data, pcm)
	if err != nil {
		t.Fatalf("Couldn't decode data: %v", err)
	}
	if len(pcm) != n {
		t.Fatalf("Length mismatch: %d samples in, %d out", len(pcm), n)
	}
}

func TestStereo(t *testing.T) {
	// Create bogus input sound
	const G4 = 391.995
	const E3 = 164.814
	const SAMPLE_RATE = 48000
	const FRAME_SIZE_MS = 60
	const CHANNELS = 2
	const FRAME_SIZE_MONO = SAMPLE_RATE * FRAME_SIZE_MS / 1000

	enc, err := NewEncoder(SAMPLE_RATE, CHANNELS, AppVoIP)
	if err != nil || enc == nil {
		t.Fatalf("Error creating new encoder: %v", err)
	}
	dec, err := NewDecoder(SAMPLE_RATE, CHANNELS)
	if err != nil || dec == nil {
		t.Fatalf("Error creating new decoder: %v", err)
	}

	// Source signal (left G4, right E3)
	left := make([]int16, FRAME_SIZE_MONO)
	right := make([]int16, FRAME_SIZE_MONO)
	addSine(left, SAMPLE_RATE, G4)
	addSine(right, SAMPLE_RATE, E3)
	pcm := interleave(left, right)

	data := make([]byte, 1000)
	n, err := enc.Encode(pcm, data)
	if err != nil {
		t.Fatalf("Couldn't encode data: %v", err)
	}
	data = data[:n]
	n, err = dec.Decode(data, pcm)
	if err != nil {
		t.Fatal("Couldn't decode data: %v", err)
	}
	if n*CHANNELS != len(pcm) {
		t.Fatalf("Length mismatch: %d samples in, %d out", len(pcm), n*CHANNELS)
	}

	// This is hard to check programatically, but I plotted the graphs in a
	// spreadsheet and it looked great. The encoded waves both start out with a
	// string of zero samples before they pick up, and the G4 is phase shifted
	// by half a period, but that's all fine for lossy audio encoding.
	/*
		leftdec, rightdec := split(pcm)
		fmt.Printf("left_in,left_out,right_in,right_out\n")
		for i := range left {
			fmt.Printf("%d,%d,%d,%d\n", left[i], leftdec[i], right[i], rightdec[i])
		}
	*/
}
