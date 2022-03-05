// Copyright © 2015-2017 Go Opus Authors (see AUTHORS file)
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
		buf[i] += int16(math.Sin(float64(i)*factor) * (math.MaxInt16 - 1))
	}
}

func maxDiff(a []int16, b []int16) int32 {
	if len(a) != len(b) {
		return math.MaxInt16
	}
	var max int32 = 0
	for i := range a {
		d := int32(a[i]) - int32(b[i])
		if d < 0 {
			d = -d
		}
		if d > max {
			max = d
		}
	}
	return max
}

func interleave(a []int16, b []int16) []int16 {
	if len(a) != len(b) {
		panic("interleave: buffers must have equal length")
	}
	result := make([]int16, 2*len(a))
	for i := range a {
		result[2*i] = a[i]
		result[2*i+1] = b[i]
	}
	return result
}

func split(interleaved []int16) ([]int16, []int16) {
	if len(interleaved)%2 != 0 {
		panic("split: interleaved buffer must have even number of samples")
	}
	left := make([]int16, len(interleaved)/2)
	right := make([]int16, len(interleaved)/2)
	for i := range left {
		left[i] = interleaved[2*i]
		right[i] = interleaved[2*i+1]
	}
	return left, right
}
