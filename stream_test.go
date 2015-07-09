// Copyright © 2015 Hraban Luyat <hraban@0brg.net>
//
// License for use of this code is detailed in the LICENSE file

package opus

import (
	"io"
	"strings"
	"testing"
)

func TestStream(t *testing.T) {
	// Simple testing, no actual decoding
	reader := strings.NewReader("hello")
	_, err := NewStream(reader)
	if err == nil {
		t.Fatalf("Expected error while initializing illegal opus stream")
	}
}

func TestCloser(t *testing.T) {
	/* TODO: test if stream.Close() also closes the underlying reader */
}
