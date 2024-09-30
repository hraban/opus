// Copyright © Go Opus Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

package opus

import (
	"fmt"
	"testing"
)

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

func TestDecoder_GetLastPacketDuration(t *testing.T) {
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
	samples, err := dec.LastPacketDuration()
	if err != nil {
		t.Fatalf("Couldn't get last packet duration: %v", err)
	}
	if samples != n {
		t.Fatalf("Wrong duration length. Expected %d. Got %d", n, samples)
	}
}

func TestDecoder_SetComplexity(t *testing.T) {
	const SAMPLE_RATE = 48000

	dec, err := NewDecoder(SAMPLE_RATE, 1)
	if err != nil || dec == nil {
		t.Fatalf("Error creating new decoder: %v", err)
	}

	t.Run("Complexity 0 to 10", func(t *testing.T) {
		for i := 0; i <= 10; i++ {
			err = dec.SetComplexity(i)
			if err != nil {
				t.Fatalf("Expected nil got %v", err)
			}
		}
	})

	for _, tt := range []int{-1, 11, 99} {
		t.Run(fmt.Sprintf("Complexity %d", tt), func(t *testing.T) {
			err = dec.SetComplexity(tt)
			if err != ErrBadArg {
				t.Fatalf("Expected %v got %v", ErrBadArg, err)
			}
		})
	}
}
