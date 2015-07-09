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

func readStreamWav(stream *Stream) ([]byte, error) {
	var buf bytes.Buffer
	pcm := make([]int16, 1000)
	for {
		n, err := stream.Read(pcm)
		switch err {
		case io.EOF:
			return buf.Bytes(), nil
		case nil:
			break
		default:
			return nil, err
		}
		if n == 0 {
			return nil, fmt.Errorf("Nil-error Read() must not return 0")
		}
		for i := 0; i < n; i++ {
			buf.WriteByte(byte(pcm[i] & 0xff))
			buf.WriteByte(byte(pcm[i] >> 8))
		}
	}
}

func TestStream(t *testing.T) {
	// Simple testing, no actual decoding
	reader, err := os.Open("testdata/speech_8.opus")
	if err != nil {
		t.Fatalf("Error while opening file: %v", err)
	}
	stream, err := NewStream(reader)
	if err != nil {
		t.Fatalf("Error while creating opus stream: %v", err)
	}
	wav, err := readStreamWav(stream)
	if err != nil {
		t.Fatalf("Error while decoding opus file: %v", err)
	}
	if len(wav) != 1036800 {
		t.Fatalf("Unexpected length of WAV file: %d", len(wav))
	}
}

func TestCloser(t *testing.T) {
	/* TODO: test if stream.Close() also closes the underlying reader */
}
