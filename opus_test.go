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

func TestCodecFEC(t *testing.T) {
	// Create bogus input sound
	const G4 = 391.995
	const SAMPLE_RATE = 48000
	const FRAME_SIZE_MS = 10
	const FRAME_SIZE = SAMPLE_RATE * FRAME_SIZE_MS / 1000
	const NUMBER_OF_FRAMES = 6
	pcm := make([]int16, 0)

	enc, err := NewEncoder(SAMPLE_RATE, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Fatalf("Error creating new encoder: %v", err)
	}
	enc.SetPacketLossPerc(30)
	enc.SetInBandFEC(true)

	dec, err := NewDecoder(SAMPLE_RATE, 1)
	if err != nil || dec == nil {
		t.Fatalf("Error creating new decoder: %v", err)
	}

	mono := make([]int16, FRAME_SIZE*NUMBER_OF_FRAMES)
	addSine(mono, SAMPLE_RATE, G4)

	encodedData := make([][]byte, NUMBER_OF_FRAMES)
	for i, idx := 0, 0; i < len(mono); i += FRAME_SIZE {
		data := make([]byte, 1000)
		n, err := enc.Encode(mono[i:i+FRAME_SIZE], data)
		if err != nil {
			t.Fatalf("Couldn't encode data: %v", err)
		}
		data = data[:n]
		encodedData[idx] = data
		idx++
	}

	lost := false
	for i := 0; i < len(encodedData); i++ {
		// "Dropping" packets 2 and 4
		if i == 2 || i == 4 {
			lost = true
			continue
		}
		if lost {
			samples, err := dec.LastPacketDuration()
			if err != nil {
				t.Fatalf("Couldn't get last packet duration: %v", err)
			}

			pcmBuffer := make([]int16, samples)
			if err = dec.DecodeFEC(encodedData[i], pcmBuffer); err != nil {
				t.Fatalf("Couldn't decode data: %v", err)
			}
			pcm = append(pcm, pcmBuffer...)
			lost = false
		}

		pcmBuffer := make([]int16, NUMBER_OF_FRAMES*FRAME_SIZE)
		n, err := dec.Decode(encodedData[i], pcmBuffer)
		if err != nil {
			t.Fatalf("Couldn't decode data: %v", err)
		}
		pcmBuffer = pcmBuffer[:n]
		if n != FRAME_SIZE {
			t.Fatalf("Length mismatch: %d samples expected, %d out", FRAME_SIZE, n)
		}
		pcm = append(pcm, pcmBuffer...)
	}

	if len(mono) != len(pcm) {
		t.Fatalf("Input/Output length mismatch: %d samples in, %d out", len(mono), len(pcm))
	}

	// Commented out for the same reason as in TestStereo
	/*
		fmt.Printf("in,out\n")
		for i := range mono {
			fmt.Printf("%d,%d\n", mono[i], pcm[i])
		}
	*/

}

func TestCodecFECFloat32(t *testing.T) {
	// Create bogus input sound
	const G4 = 391.995
	const SAMPLE_RATE = 48000
	const FRAME_SIZE_MS = 10
	const FRAME_SIZE = SAMPLE_RATE * FRAME_SIZE_MS / 1000
	const NUMBER_OF_FRAMES = 6
	pcm := make([]float32, 0)

	enc, err := NewEncoder(SAMPLE_RATE, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Fatalf("Error creating new encoder: %v", err)
	}
	enc.SetPacketLossPerc(30)
	enc.SetInBandFEC(true)

	dec, err := NewDecoder(SAMPLE_RATE, 1)
	if err != nil || dec == nil {
		t.Fatalf("Error creating new decoder: %v", err)
	}

	mono := make([]float32, FRAME_SIZE*NUMBER_OF_FRAMES)
	addSineFloat32(mono, SAMPLE_RATE, G4)

	encodedData := make([][]byte, NUMBER_OF_FRAMES)
	for i, idx := 0, 0; i < len(mono); i += FRAME_SIZE {
		data := make([]byte, 1000)
		n, err := enc.EncodeFloat32(mono[i:i+FRAME_SIZE], data)
		if err != nil {
			t.Fatalf("Couldn't encode data: %v", err)
		}
		data = data[:n]
		encodedData[idx] = data
		idx++
	}

	lost := false
	for i := 0; i < len(encodedData); i++ {
		// "Dropping" packets 2 and 4
		if i == 2 || i == 4 {
			lost = true
			continue
		}
		if lost {
			samples, err := dec.LastPacketDuration()
			if err != nil {
				t.Fatalf("Couldn't get last packet duration: %v", err)
			}

			pcmBuffer := make([]float32, samples)
			if err = dec.DecodeFECFloat32(encodedData[i], pcmBuffer); err != nil {
				t.Fatalf("Couldn't decode data: %v", err)
			}
			pcm = append(pcm, pcmBuffer...)
			lost = false
		}

		pcmBuffer := make([]float32, NUMBER_OF_FRAMES*FRAME_SIZE)
		n, err := dec.DecodeFloat32(encodedData[i], pcmBuffer)
		if err != nil {
			t.Fatalf("Couldn't decode data: %v", err)
		}
		pcmBuffer = pcmBuffer[:n]
		if n != FRAME_SIZE {
			t.Fatalf("Length mismatch: %d samples expected, %d out", FRAME_SIZE, n)
		}
		pcm = append(pcm, pcmBuffer...)
	}

	if len(mono) != len(pcm) {
		t.Fatalf("Input/Output length mismatch: %d samples in, %d out", len(mono), len(pcm))
	}

	// Commented out for the same reason as in TestStereo
	/*
		fmt.Printf("in,out\n")
		for i := 0; i < len(mono); i++ {
			fmt.Printf("%f,%f\n", mono[i], pcm[i])
		}
	*/
}

func TestCodecPLC(t *testing.T) {
	// Create bogus input sound
	const G4 = 391.995
	const SAMPLE_RATE = 48000
	const FRAME_SIZE_MS = 10
	const FRAME_SIZE = SAMPLE_RATE * FRAME_SIZE_MS / 1000
	const NUMBER_OF_FRAMES = 6
	pcm := make([]int16, 0)

	enc, err := NewEncoder(SAMPLE_RATE, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Fatalf("Error creating new encoder: %v", err)
	}

	dec, err := NewDecoder(SAMPLE_RATE, 1)
	if err != nil || dec == nil {
		t.Fatalf("Error creating new decoder: %v", err)
	}

	mono := make([]int16, FRAME_SIZE*NUMBER_OF_FRAMES)
	addSine(mono, SAMPLE_RATE, G4)

	encodedData := make([][]byte, NUMBER_OF_FRAMES)
	for i, idx := 0, 0; i < len(mono); i += FRAME_SIZE {
		data := make([]byte, 1000)
		n, err := enc.Encode(mono[i:i+FRAME_SIZE], data)
		if err != nil {
			t.Fatalf("Couldn't encode data: %v", err)
		}
		data = data[:n]
		encodedData[idx] = data
		idx++
	}

	lost := false
	for i := 0; i < len(encodedData); i++ {
		// "Dropping" packets 2 and 4
		if i == 2 || i == 4 {
			lost = true
			continue
		}
		if lost {
			samples, err := dec.LastPacketDuration()
			if err != nil {
				t.Fatalf("Couldn't get last packet duration: %v", err)
			}

			pcmBuffer := make([]int16, samples)
			if err = dec.DecodePLC(pcmBuffer); err != nil {
				t.Fatalf("Couldn't decode data: %v", err)
			}
			nonZero := 0
			for n := range pcmBuffer {
				if pcmBuffer[n] != 0 {
					nonZero++
				}
			}
			if nonZero == 0 {
				t.Fatalf("DecodePLC produced only zero samples")
			}
			pcm = append(pcm, pcmBuffer...)
			lost = false
		}

		pcmBuffer := make([]int16, NUMBER_OF_FRAMES*FRAME_SIZE)
		n, err := dec.Decode(encodedData[i], pcmBuffer)
		if err != nil {
			t.Fatalf("Couldn't decode data: %v", err)
		}
		pcmBuffer = pcmBuffer[:n]
		if n != FRAME_SIZE {
			t.Fatalf("Length mismatch: %d samples expected, %d out", FRAME_SIZE, n)
		}
		pcm = append(pcm, pcmBuffer...)
	}

	if len(mono) != len(pcm) {
		t.Fatalf("Input/Output length mismatch: %d samples in, %d out", len(mono), len(pcm))
	}

	// Commented out for the same reason as in TestStereo
	/*
		fmt.Printf("in,out\n")
		for i := range mono {
			fmt.Printf("%d,%d\n", mono[i], pcm[i])
		}
	*/

}

func TestCodecPLCFloat32(t *testing.T) {
	// Create bogus input sound
	const G4 = 391.995
	const SAMPLE_RATE = 48000
	const FRAME_SIZE_MS = 10
	const FRAME_SIZE = SAMPLE_RATE * FRAME_SIZE_MS / 1000
	const NUMBER_OF_FRAMES = 6
	pcm := make([]float32, 0)

	enc, err := NewEncoder(SAMPLE_RATE, 1, AppVoIP)
	if err != nil || enc == nil {
		t.Fatalf("Error creating new encoder: %v", err)
	}

	dec, err := NewDecoder(SAMPLE_RATE, 1)
	if err != nil || dec == nil {
		t.Fatalf("Error creating new decoder: %v", err)
	}

	mono := make([]float32, FRAME_SIZE*NUMBER_OF_FRAMES)
	addSineFloat32(mono, SAMPLE_RATE, G4)

	encodedData := make([][]byte, NUMBER_OF_FRAMES)
	for i, idx := 0, 0; i < len(mono); i += FRAME_SIZE {
		data := make([]byte, 1000)
		n, err := enc.EncodeFloat32(mono[i:i+FRAME_SIZE], data)
		if err != nil {
			t.Fatalf("Couldn't encode data: %v", err)
		}
		data = data[:n]
		encodedData[idx] = data
		idx++
	}

	lost := false
	for i := 0; i < len(encodedData); i++ {
		// "Dropping" packets 2 and 4
		if i == 2 || i == 4 {
			lost = true
			continue
		}
		if lost {
			samples, err := dec.LastPacketDuration()
			if err != nil {
				t.Fatalf("Couldn't get last packet duration: %v", err)
			}

			pcmBuffer := make([]float32, samples)
			if err = dec.DecodePLCFloat32(pcmBuffer); err != nil {
				t.Fatalf("Couldn't decode data: %v", err)
			}
			nonZero := 0
			for n := range pcmBuffer {
				if pcmBuffer[n] != 0.0 {
					nonZero++
				}
			}
			if nonZero == 0 {
				t.Fatalf("DecodePLC produced only zero samples")
			}
			pcm = append(pcm, pcmBuffer...)
			lost = false
		}

		pcmBuffer := make([]float32, NUMBER_OF_FRAMES*FRAME_SIZE)
		n, err := dec.DecodeFloat32(encodedData[i], pcmBuffer)
		if err != nil {
			t.Fatalf("Couldn't decode data: %v", err)
		}
		pcmBuffer = pcmBuffer[:n]
		if n != FRAME_SIZE {
			t.Fatalf("Length mismatch: %d samples expected, %d out", FRAME_SIZE, n)
		}
		pcm = append(pcm, pcmBuffer...)
	}

	if len(mono) != len(pcm) {
		t.Fatalf("Input/Output length mismatch: %d samples in, %d out", len(mono), len(pcm))
	}

	// Commented out for the same reason as in TestStereo
	/*
		fmt.Printf("in,out\n")
		for i := 0; i < len(mono); i++ {
			fmt.Printf("%f,%f\n", mono[i], pcm[i])
		}
	*/
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
		t.Fatalf("Couldn't decode data: %v", err)
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

func TestBufferSize(t *testing.T) {
	const G4 = 391.995
	const SAMPLE_RATE = 48000
	const FRAME_SIZE_MS = 60
	const FRAME_SIZE = SAMPLE_RATE * FRAME_SIZE_MS / 1000
	const GUARD_SIZE = 100

	checkGuardInt16 := func(t *testing.T, s []int16) {
		for n := range s {
			if s[n] != 0 {
				t.Fatal("Memory corruption detected")
			}
		}
	}

	checkGuardFloat32 := func(t *testing.T, s []float32) {
		for n := range s {
			if s[n] != 0 {
				t.Fatal("Memory corruption detected")
			}
		}
	}

	checkResult := func(t *testing.T, n int, err error, expectOK bool) {
		if expectOK {
			if err != nil {
				t.Fatalf("Couldn't decode data: %v", err)
			}
			if n != FRAME_SIZE {
				t.Fatalf("Length mismatch: %d samples in, %d out", FRAME_SIZE, n)
			}
		} else {
			if err == nil {
				t.Fatalf("Expected decode failure, but it succeeded")
			}
		}
	}

	encodeFrame := func(t *testing.T) []byte {
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
		return data[:n]
	}

	createDecoder := func(t *testing.T) *Decoder {
		dec, err := NewDecoder(SAMPLE_RATE, 1)
		if err != nil || dec == nil {
			t.Fatalf("Error creating new decoder: %v", err)
		}
		return dec
	}

	decodeInt16 := func(t *testing.T, data []byte, decodeSize int, expectOK bool) {
		dec := createDecoder(t)
		decodedMem := make([]int16, decodeSize+GUARD_SIZE*2)
		decodedRef := decodedMem[GUARD_SIZE : GUARD_SIZE+decodeSize : GUARD_SIZE+decodeSize]
		n, err := dec.Decode(data, decodedRef)
		checkGuardInt16(t, decodedMem[:GUARD_SIZE])
		checkGuardInt16(t, decodedMem[decodeSize+GUARD_SIZE:])
		checkResult(t, n, err, expectOK)
	}

	decodeFloat32 := func(t *testing.T, data []byte, decodeSize int, expectOK bool) {
		dec := createDecoder(t)
		decodedMem := make([]float32, decodeSize+GUARD_SIZE*2)
		decodedRef := decodedMem[GUARD_SIZE : GUARD_SIZE+decodeSize : GUARD_SIZE+decodeSize]
		n, err := dec.DecodeFloat32(data, decodedRef)
		checkGuardFloat32(t, decodedMem[:GUARD_SIZE])
		checkGuardFloat32(t, decodedMem[decodeSize+GUARD_SIZE:])
		checkResult(t, n, err, expectOK)
	}

	decodeFecInt16 := func(t *testing.T, data []byte, decodeSize int, expectOK bool) {
		dec := createDecoder(t)
		decodedMem := make([]int16, decodeSize+GUARD_SIZE*2)
		decodedRef := decodedMem[GUARD_SIZE : GUARD_SIZE+decodeSize : GUARD_SIZE+decodeSize]
		err := dec.DecodeFEC(data, decodedRef)
		checkGuardInt16(t, decodedMem[:GUARD_SIZE])
		checkGuardInt16(t, decodedMem[decodeSize+GUARD_SIZE:])
		checkResult(t, FRAME_SIZE, err, expectOK)
	}

	decodeFecFloat32 := func(t *testing.T, data []byte, decodeSize int, expectOK bool) {
		dec := createDecoder(t)
		decodedMem := make([]float32, decodeSize+GUARD_SIZE*2)
		decodedRef := decodedMem[GUARD_SIZE : GUARD_SIZE+decodeSize : GUARD_SIZE+decodeSize]
		err := dec.DecodeFECFloat32(data, decodedRef)
		checkGuardFloat32(t, decodedMem[:GUARD_SIZE])
		checkGuardFloat32(t, decodedMem[decodeSize+GUARD_SIZE:])
		checkResult(t, FRAME_SIZE, err, expectOK)
	}

	t.Run("smaller-buffer-int16", func(t *testing.T) {
		decodeInt16(t, encodeFrame(t), FRAME_SIZE-1, false)
	})

	t.Run("smaller-buffer-float32", func(t *testing.T) {
		decodeFloat32(t, encodeFrame(t), FRAME_SIZE-1, false)
	})

	t.Run("smaller-buffer-int16-fec", func(t *testing.T) {
		decodeFecFloat32(t, encodeFrame(t), FRAME_SIZE-1, false)
	})

	t.Run("smaller-buffer-float32-fec", func(t *testing.T) {
		decodeFecFloat32(t, encodeFrame(t), FRAME_SIZE-1, false)
	})

	t.Run("exact-buffer-int16", func(t *testing.T) {
		decodeInt16(t, encodeFrame(t), FRAME_SIZE, true)
	})

	t.Run("exact-buffer-float32", func(t *testing.T) {
		decodeFloat32(t, encodeFrame(t), FRAME_SIZE, true)
	})

	t.Run("exact-buffer-int16-fec", func(t *testing.T) {
		decodeFecInt16(t, encodeFrame(t), FRAME_SIZE, true)
	})

	t.Run("exact-buffer-float32-fec", func(t *testing.T) {
		decodeFecFloat32(t, encodeFrame(t), FRAME_SIZE, true)
	})

	t.Run("larger-buffer-int16", func(t *testing.T) {
		decodeInt16(t, encodeFrame(t), FRAME_SIZE+1, true)
	})

	t.Run("larger-buffer-float32", func(t *testing.T) {
		decodeFloat32(t, encodeFrame(t), FRAME_SIZE+1, true)
	})

	t.Run("larger-buffer-int16-fec", func(t *testing.T) {
		decodeFecInt16(t, encodeFrame(t), FRAME_SIZE+1, false)
	})

	t.Run("larger-buffer-float32-fec", func(t *testing.T) {
		decodeFecFloat32(t, encodeFrame(t), FRAME_SIZE+1, false)
	})
}
