// Copyright © 2015-2017 Go Opus Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

package opus

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestStreamIllegal(t *testing.T) {
	// Simple testing, no actual decoding
	reader := strings.NewReader("hello test test this is not a legal Opus stream")
	_, err := NewStream(reader)
	if err == nil {
		t.Fatalf("Expected error while initializing illegal opus stream")
	}
}

func readStreamPcm(t *testing.T, stream *Stream, buffersize int) []int16 {
	var pcm []int16
	pcmbuf := make([]int16, buffersize)
	for {
		n, err := stream.Read(pcmbuf)
		switch err {
		case io.EOF:
			return pcm
		case nil:
			break
		default:
			t.Fatalf("Error while decoding opus file: %v", err)
		}
		if n == 0 {
			t.Fatal("Nil-error Read() must not return 0")
		}
		pcm = append(pcm, pcmbuf[:n]...)
	}
}

func mustOpenFile(t *testing.T, fname string) *os.File {
	f, err := os.Open(fname)
	if err != nil {
		t.Fatalf("Error while opening %s: %v", fname, err)
		return nil
	}
	return f
}

func mustOpenStream(t *testing.T, r io.Reader) *Stream {
	stream, err := NewStream(r)
	if err != nil {
		t.Fatalf("Error while creating opus stream: %v", err)
		return nil
	}
	return stream
}

func opus2pcm(t *testing.T, fname string, buffersize int) []int16 {
	reader := mustOpenFile(t, fname)
	stream := mustOpenStream(t, reader)
	return readStreamPcm(t, stream, buffersize)
}

// Extract raw pcm data from .wav file
func extractWavPcm(t *testing.T, fname string) []int16 {
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("Error reading file data from %s: %v", fname, err)
	}
	const wavHeaderSize = 44
	if (len(bytes)-wavHeaderSize)%2 == 1 {
		t.Fatalf("Illegal wav data: payload must be encoded in byte pairs")
	}
	numSamples := (len(bytes) - wavHeaderSize) / 2
	samples := make([]int16, numSamples)
	for i := 0; i < numSamples; i++ {
		samples[i] += int16(bytes[wavHeaderSize+i*2])
		samples[i] += int16(bytes[wavHeaderSize+i*2+1]) << 8
	}
	return samples
}

func TestStream(t *testing.T) {
	opuspcm := opus2pcm(t, "testdata/speech_8.opus", 10000)
	wavpcm := extractWavPcm(t, "testdata/speech_8.wav")
	if len(opuspcm) != len(wavpcm) {
		t.Fatalf("Unexpected length of decoded opus file: %d (.wav: %d)", len(opuspcm), len(wavpcm))
	}
	d := maxDiff(opuspcm, wavpcm)
	// No science behind this number
	const epsilon = 18
	if d > epsilon {
		t.Errorf("Maximum difference between decoded streams too high: %d", d)
	}
}

func TestStreamSmallBuffer(t *testing.T) {
	smallbuf := opus2pcm(t, "testdata/speech_8.opus", 1)
	bigbuf := opus2pcm(t, "testdata/speech_8.opus", 10000)
	if !reflect.DeepEqual(smallbuf, bigbuf) {
		t.Errorf("Reading with 1-sample buffer size yields different audio data")
	}
}

type mockCloser struct {
	io.Reader
	closed bool
}

func (m *mockCloser) Close() error {
	if m.closed {
		return fmt.Errorf("Already closed")
	}
	m.closed = true
	return nil
}

func TestCloser(t *testing.T) {
	f := mustOpenFile(t, "testdata/speech_8.opus")
	mc := &mockCloser{Reader: f}
	stream := mustOpenStream(t, mc)
	stream.Close()
	if !mc.closed {
		t.Error("Expected opus stream to call .Close on the reader")
	}
}
