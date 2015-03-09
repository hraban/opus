// Copyright © 2015 Hraban Luyat <hraban@0brg.net>
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

func addSine(buf []float32, sampleRate int, freq float64) {
	factor := 2 * math.Pi * freq / float64(sampleRate)
	for i := range buf {
		buf[i] += float32(math.Sin(float64(i) * factor))
	}
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
	addSine(pcm, SAMPLE_RATE, G4)
	data, err := enc.EncodeFloat32(pcm)
	if err != nil {
		t.Fatalf("Couldn't encode data: %v", err)
	}
	dec, err := NewDecoder(SAMPLE_RATE, 1)
	if err != nil || dec == nil {
		t.Fatalf("Error creating new decoder: %v", err)
	}
	pcm2, err := dec.DecodeFloat32(data)
	if err != nil {
		t.Fatalf("Couldn't decode data: %v", err)
	}
	if len(pcm) != len(pcm2) {
		t.Fatalf("Length mismatch: %d samples in, %d out", len(pcm), len(pcm2))
	}
	// Checking the output programmatically is seriously not easy. I checked it
	// by hand and by ear, it looks fine. I'll just assume the API faithfully
	// passes error codes through, and that's that.
}
