// Copyright © 2015-2017 Go Opus Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

package opus

import (
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
