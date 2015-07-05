// Copyright © 2015 Hraban Luyat <hraban@0brg.net>
//
// License for use of this code is detailed in the LICENSE file

package opus

import (
	"strings"
	"testing"
)

func Test(t *testing.T) {
	// Simple testing, no actual decoding
	reader := strings.NewReader("hello")
	_, err := NewStream(reader)
	if err == nil {
		t.Fatalf("Expected error while initializing illegal opus stream")
	}
}
