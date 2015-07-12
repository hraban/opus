// Copyright © 2015 Hraban Luyat <hraban@0brg.net>
//
// License for use of this code is detailed in the LICENSE file

package opus

import (
	"bytes"
	"fmt"
	"io"
	"os"
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

func readStreamWav(t *testing.T, stream *Stream, buffersize int) []byte {
	var buf bytes.Buffer
	pcm := make([]int16, buffersize)
	for {
		n, err := stream.Read(pcm)
		switch err {
		case io.EOF:
			return buf.Bytes()
		case nil:
			break
		default:
			t.Fatalf("Error while decoding opus file: %v", err)
		}
		if n == 0 {
			t.Fatal("Nil-error Read() must not return 0")
		}
		for i := 0; i < n; i++ {
			buf.WriteByte(byte(pcm[i] & 0xff))
			buf.WriteByte(byte(pcm[i] >> 8))
		}
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

func readFileWav(t *testing.T, fname string, buffersize int) []byte {
	reader := mustOpenFile(t, fname)
	stream, err := NewStream(reader)
	if err != nil {
		t.Fatalf("Error while creating opus stream: %v", err)
	}
	return readStreamWav(t, stream, buffersize)
}

func TestStream(t *testing.T) {
	wav := readFileWav(t, "testdata/speech_8.opus", 10000)
	if len(wav) != 1036800 {
		t.Fatalf("Unexpected length of WAV file: %d", len(wav))
	}
}

func TestStreamSmallBuffer(t *testing.T) {
	smallbuf := readFileWav(t, "testdata/speech_8.opus", 1)
	bigbuf := readFileWav(t, "testdata/speech_8.opus", 10000)
	if !bytes.Equal(smallbuf, bigbuf) {
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
