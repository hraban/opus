// Copyright © 2015, 2016 Hraban Luyat <hraban@0brg.net>
//
// License for use of this code is detailed in the LICENSE file

package opus

import (
	"math"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	if ver := Version(); !strings.HasPrefix(ver, "libopus") {
		t.Errorf("Unexpected linked libopus version: " + ver)
	}
}

func TestEncoderNew(t *testing.T) {
	enc, err := NewEncoder(48000, 1, APPLICATION_VOIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	enc, err = NewEncoder(12345, 1, APPLICATION_VOIP)
	if err == nil || enc != nil {
		t.Errorf("Expected error for illegal samplerate 12345")
	}
}

func TestDecoderNew(t *testing.T) {
	dec, err := NewDecoder(48000, 1)
	if err != nil || dec == nil {
		t.Errorf("Error creating new decoder: %v", err)
	}
	dec, err = NewDecoder(12345, 1)
	if err == nil || dec != nil {
		t.Errorf("Expected error for illegal samplerate 12345")
	}
}

func TestEncoderUnitialized(t *testing.T) {
	var enc Encoder
	_, err := enc.Encode(nil, nil)
	if err != errEncUninitialized {
		t.Errorf("Expected \"unitialized encoder\" error: %v", err)
	}
	_, err = enc.EncodeFloat32(nil, nil)
	if err != errEncUninitialized {
		t.Errorf("Expected \"unitialized encoder\" error: %v", err)
	}
}

func TestDecoderUnitialized(t *testing.T) {
	var dec Decoder
	_, err := dec.Decode(nil, nil)
	if err != errDecUninitialized {
		t.Errorf("Expected \"unitialized decoder\" error: %v", err)
	}
	_, err = dec.DecodeFloat32(nil, nil)
	if err != errDecUninitialized {
		t.Errorf("Expected \"unitialized decoder\" error: %v", err)
	}
}

func TestOpusErrstr(t *testing.T) {
	// I scooped this -1 up from opus_defines.h, it's OPUS_BAD_ARG. Not pretty,
	// but it's better than not testing at all. Again, accessing #defines from
	// CGO is not possible.
	if ERR_OPUS_BAD_ARG.Error() != "opus: invalid argument" {
		t.Errorf("Expected \"invalid argument\" error message for error code %d: %v",
			ERR_OPUS_BAD_ARG, ERR_OPUS_BAD_ARG)
	}
}

func addSineFloat32(buf []float32, sampleRate int, freq float64) {
	factor := 2 * math.Pi * freq / float64(sampleRate)
	for i := range buf {
		buf[i] += float32(math.Sin(float64(i) * factor))
	}
}

func addSine(buf []int16, sampleRate int, freq float64) {
	factor := 2 * math.Pi * freq / float64(sampleRate)
	for i := range buf {
		buf[i] += int16(math.Sin(float64(i)*factor) * math.MaxInt16)
	}
}

func TestCodec(t *testing.T) {
	// Create bogus input sound
	const G4 = 391.995
	const SAMPLE_RATE = 48000
	const FRAME_SIZE_MS = 60
	const FRAME_SIZE = SAMPLE_RATE * FRAME_SIZE_MS / 1000
	pcm := make([]int16, FRAME_SIZE)
	enc, err := NewEncoder(SAMPLE_RATE, 1, APPLICATION_VOIP)
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
	enc, err := NewEncoder(SAMPLE_RATE, 1, APPLICATION_VOIP)
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
	n, err = dec.DecodeFloat32(data, pcm)
	if err != nil {
		t.Fatalf("Couldn't decode data: %v", err)
	}
	if len(pcm) != n {
		t.Fatalf("Length mismatch: %d samples in, %d out", len(pcm), n)
	}
}
