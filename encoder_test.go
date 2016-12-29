package opus

import (
	"testing"
)

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

func TestEncoderDTX(t *testing.T) {
	enc, err := NewEncoder(8000, 1, APPLICATION_VOIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	vals := []bool{true, false}
	for _, dtx := range vals {
		enc.UseDTX(dtx)
		gotv := enc.DTX()
		if gotv != dtx {
			t.Errorf("Error set dtx: expect dtx=%v, got dtx=%v", dtx, gotv)
		}
	}
}

func TestEncoderSampleRate(t *testing.T) {
	sample_rates := []int{8000, 12000, 16000, 24000, 48000}
	for _, f := range sample_rates {
		enc, err := NewEncoder(f, 1, APPLICATION_VOIP)
		if err != nil || enc == nil {
			t.Fatalf("Error creating new encoder with sample_rate %d Hz: %v", f, err)
		}
		f2 := enc.SampleRate()
		if f != f2 {
			t.Errorf("Unexpected sample rate reported by %d Hz encoder: %d", f, f2)
		}
	}
}

func TestEncoder_SetGetBitrate(t *testing.T) {
	enc, err := NewEncoder(8000, 1, APPLICATION_AUDIO)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	vals := []int{500, 100000, 300000}
	for _, bitrate := range vals {
		err := enc.SetBitrate(bitrate)
		if err != nil {
			t.Error("Error set bitrate:", err)
		}
		br, err := enc.Bitrate()
		if err != nil {
			t.Error("Error getting bitrate", err)
		}
		if br != bitrate {
			t.Errorf("Unexpected bitrate. Got %d, but expected %d", br, bitrate)
		}
	}
}

func TestEncoder_SetGetInvalidBitrate(t *testing.T) {
	enc, err := NewEncoder(8000, 1, APPLICATION_AUDIO)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	invalidVals := []int{-20, 0}
	for _, bitrate := range invalidVals {
		err := enc.SetBitrate(bitrate)
		if err == nil {
			t.Errorf("Expected Error invalid bitrate: %d", bitrate)
		}
		br, err := enc.Bitrate()
		if err != nil {
			t.Error("Error getting bitrate", err)
		}
		// default bitrate is 32 kbit/s
		if br != 32000 {
			t.Errorf("Unexpected bitrate. Got %d, but expected %d", br, bitrate)
		}
	}
}

func TestEncoder_SetGetComplexity(t *testing.T) {
	enc, err := NewEncoder(8000, 1, APPLICATION_AUDIO)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	vals := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, complexity := range vals {
		err := enc.SetComplexity(complexity)
		if err != nil {
			t.Error("Error setting complexity value:", err)
		}
		cpx, err := enc.Complexity()
		if err != nil {
			t.Error("Error getting complexity value", err)
		}
		if cpx != complexity {
			t.Errorf("Unexpected encoder complexity value. Got %d, but expected %d",
				cpx, complexity)
		}
	}
}

func TestEncoder_SetGetInvalidComplexity(t *testing.T) {
	enc, err := NewEncoder(8000, 1, APPLICATION_AUDIO)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	invalidVals := []int{-20, 11, 1000}
	for _, complexity := range invalidVals {
		err := enc.SetComplexity(complexity)
		if err == nil {
			t.Errorf("Expected Error invalid complexity value: %d", complexity)
		}
		cpx, err := enc.Complexity()
		if err != nil {
			t.Error("Error getting complexity value", err)
		}
		// default complexity value is 9
		if cpx != 9 {
			t.Errorf("Unexpected complexity value. Got %d, but expected %d",
				cpx, complexity)
		}
	}
}

func TestEncoder_SetGetMaxBandwidth(t *testing.T) {
	enc, err := NewEncoder(8000, 1, APPLICATION_AUDIO)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	vals := []int{OPUS_BANDWIDTH_NARROWBAND,
		OPUS_BANDWIDTH_MEDIUMBAND,
		OPUS_BANDWIDTH_WIDEBAND,
		OPUS_BANDWIDTH_SUPERWIDEBAND,
		OPUS_BANDWIDTH_FULLBAND,
	}
	for _, maxBw := range vals {
		err := enc.SetMaxBandwidth(maxBw)
		if err != nil {
			t.Error("Error setting max Bandwidth:", err)
		}
		maxBwRead, err := enc.MaxBandwidth()
		if err != nil {
			t.Error("Error getting max Bandwidth", err)
		}
		if maxBwRead != maxBw {
			t.Errorf("Unexpected max Bandwidth value. Got %d, but expected %d",
				maxBwRead, maxBw)
		}
	}
}
