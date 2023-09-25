// Copyright © Go Opus Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

package opus

import "testing"

func TestEncoderNew(t *testing.T) {
	enc, err := NewEncoder(48000, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	enc, err = NewEncoder(12345, 1, AppVoIP)
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
	enc, err := NewEncoder(8000, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	vals := []bool{true, false}
	for _, dtx := range vals {
		err := enc.SetDTX(dtx)
		if err != nil {
			t.Fatalf("Error setting DTX to %t: %v", dtx, err)
		}
		gotv, err := enc.DTX()
		if err != nil {
			t.Fatalf("Error getting DTX (%t): %v", dtx, err)
		}
		if gotv != dtx {
			t.Errorf("Error set dtx: expect dtx=%v, got dtx=%v", dtx, gotv)
		}
	}
}

func TestEncoderInDTX(t *testing.T) {
	// Create bogus input sound
	const G4 = 391.995
	const SAMPLE_RATE = 48000
	const FRAME_SIZE_MS = 60
	const FRAME_SIZE = SAMPLE_RATE * FRAME_SIZE_MS / 1000
	pcm := make([]int16, FRAME_SIZE)
	silentPCM := make([]int16, FRAME_SIZE)
	out := make([]byte, FRAME_SIZE*4)
	addSine(pcm, SAMPLE_RATE, G4)

	vals := []bool{true, false}
	for _, dtx := range vals {
		enc, err := NewEncoder(SAMPLE_RATE, 1, AppVoIP)
		if err != nil || enc == nil {
			t.Fatalf("Error creating new encoder: %v", err)
		}
		enc.SetDTX(dtx)
		if _, err = enc.Encode(pcm, out); err != nil {
			t.Fatalf("Error encoding non-silent frame: %v", err)
		}
		gotDTX, err := enc.InDTX()
		if err != nil {
			t.Fatalf("Error getting in DTX (%t): %v", dtx, err)
		}
		if gotDTX {
			t.Fatalf("Error get in dtx: expect in dtx=false, got=true")
		}
		// Encode a few frames to let DTX kick in
		for i := 0; i < 5; i++ {
			if _, err = enc.Encode(silentPCM, out); err != nil {
				t.Fatalf("Error encoding silent frame: %v", err)
			}
		}
		gotDTX, err = enc.InDTX()
		if err != nil {
			t.Fatalf("Error getting in DTX (%t): %v", dtx, err)
		}
		if gotDTX != dtx {
			t.Errorf("Error set dtx: expect in dtx=%v, got in dtx=%v", dtx, gotDTX)
		}
	}
}

func TestEncoderSampleRate(t *testing.T) {
	sample_rates := []int{8000, 12000, 16000, 24000, 48000}
	for _, f := range sample_rates {
		enc, err := NewEncoder(f, 1, AppVoIP)
		if err != nil || enc == nil {
			t.Fatalf("Error creating new encoder with sample_rate %d Hz: %v", f, err)
		}
		f2, err := enc.SampleRate()
		if err != nil {
			t.Fatalf("Error getting sample rate (%d Hz): %v", f, err)
		}
		if f != f2 {
			t.Errorf("Unexpected sample rate reported by %d Hz encoder: %d", f, f2)
		}
	}
}

func TestEncoder_SetGetBitrate(t *testing.T) {
	enc, err := NewEncoder(8000, 1, AppVoIP)
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

func TestEncoder_SetBitrateToAuto(t *testing.T) {
	enc, err := NewEncoder(8000, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}

	bitrate := 5000
	if err := enc.SetBitrate(bitrate); err != nil {
		t.Error("Error setting bitrate:", err)
	}

	br, err := enc.Bitrate()
	if err != nil {
		t.Error("Error getting bitrate", err)
	}

	if br != bitrate {
		t.Errorf("Unexpected bitrate. Got %d, but expected %d", br, bitrate)
	}

	err = enc.SetBitrateToAuto()
	if err != nil {
		t.Error("Error setting Auto bitrate:", err)
	}

	br, err = enc.Bitrate()
	if err != nil {
		t.Error("Error getting bitrate", err)
	}

	brDefault := 32000 //default start value
	if br != brDefault {
		t.Errorf("Unexpected bitrate. Got %d, but expected %d", br, brDefault)
	}
}

func TestEncoder_SetBitrateToMax(t *testing.T) {
	enc, err := NewEncoder(8000, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}

	bitrate := 5000
	if err := enc.SetBitrate(bitrate); err != nil {
		t.Error("Error setting bitrate:", err)
	}

	br, err := enc.Bitrate()
	if err != nil {
		t.Error("Error getting bitrate", err)
	}

	if br != bitrate {
		t.Errorf("Unexpected bitrate. Got %d, but expected %d", br, bitrate)
	}

	err = enc.SetBitrateToMax()
	if err != nil {
		t.Error("Error setting Max bitrate:", err)
	}

	br, err = enc.Bitrate()
	if err != nil {
		t.Error("Error getting bitrate", err)
	}

	brMax := 4083200
	if br != brMax { //default start value
		t.Errorf("Unexpected bitrate. Got %d, but expected %d", br, brMax)
	}
}

func TestEncoder_SetGetInvalidBitrate(t *testing.T) {
	enc, err := NewEncoder(8000, 1, AppVoIP)
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
	enc, err := NewEncoder(8000, 1, AppVoIP)
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
	enc, err := NewEncoder(8000, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	invalidVals := []int{-20, 11, 1000}
	for _, complexity := range invalidVals {
		err := enc.SetComplexity(complexity)
		if err == nil {
			t.Errorf("Expected Error invalid complexity value: %d", complexity)
		}
		if err.Error() != "opus: invalid argument" {
			t.Error("Unexpected Error message")
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
	enc, err := NewEncoder(8000, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	vals := []Bandwidth{
		Narrowband,
		Mediumband,
		Wideband,
		SuperWideband,
		Fullband,
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

func TestEncoder_SetGetInBandFEC(t *testing.T) {
	enc, err := NewEncoder(8000, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}

	if err := enc.SetInBandFEC(true); err != nil {
		t.Error("Error setting fec:", err)
	}

	fec, err := enc.InBandFEC()
	if err != nil {
		t.Error("Error getting fec", err)
	}
	if !fec {
		t.Errorf("Wrong fec value. Expected %t", true)
	}

	if err := enc.SetInBandFEC(false); err != nil {
		t.Error("Error setting fec:", err)
	}

	fec, err = enc.InBandFEC()
	if err != nil {
		t.Error("Error getting fec", err)
	}
	if fec {
		t.Errorf("Wrong fec value. Expected %t", false)
	}
}

func TestEncoder_SetGetPacketLossPerc(t *testing.T) {
	enc, err := NewEncoder(8000, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	vals := []int{0, 5, 10, 20}
	for _, lossPerc := range vals {
		err := enc.SetPacketLossPerc(lossPerc)
		if err != nil {
			t.Error("Error setting loss percentage value:", err)
		}
		lp, err := enc.PacketLossPerc()
		if err != nil {
			t.Error("Error getting loss percentage value", err)
		}
		if lp != lossPerc {
			t.Errorf("Unexpected encoder loss percentage value. Got %d, but expected %d",
				lp, lossPerc)
		}
	}
}

func TestEncoder_SetGetInvalidPacketLossPerc(t *testing.T) {
	enc, err := NewEncoder(8000, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	vals := []int{-1, 101}
	for _, lossPerc := range vals {
		err := enc.SetPacketLossPerc(lossPerc)
		if err == nil {
			t.Errorf("Expected Error invalid loss percentage: %d", lossPerc)
		}
		lp, err := enc.PacketLossPerc()
		if err != nil {
			t.Error("Error getting loss percentage value", err)
		}
		// default packet loss percentage is 0
		if lp != 0 {
			t.Errorf("Unexpected encoder loss percentage value. Got %d, but expected %d",
				lp, lossPerc)
		}
	}
}

func TestEncoder_Reset(t *testing.T) {
	enc, err := NewEncoder(48000, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Errorf("Error creating new encoder: %v", err)
	}
	RunTestCodec(t, enc)
	if err := enc.Reset(); err != nil {
		t.Errorf("Error reset encoder: %v", err)
	}
	RunTestCodec(t, enc)
}
