// Copyright © 2015, 2016 Hraban Luyat <hraban@0brg.net>
//
// License for use of this code is detailed in the LICENSE file

package opus

import (
	"math"
)

// utility functions for unit tests

func addSineFloat32(buf []float32, sampleRate int, freq float64) {
	factor := 2 * math.Pi * freq / float64(sampleRate)
	for i := range buf {
		buf[i] += float32(math.Sin(float64(i) * factor))
	}
}

func addSine(buf []int16, sampleRate int, freq float64) {
	factor := 2 * math.Pi * freq / float64(sampleRate)
	for i := range buf {
		buf[i] += int16(math.Sin(float64(i)*factor) * math.MaxInt16)
	}
}
