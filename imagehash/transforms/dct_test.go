// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transforms

import (
	"math/rand"
	"testing"
)

const (
	EPSILON float64 = 0.00000001
)

func TestDCT1D(t *testing.T) {
	for _, tt := range []struct {
		input  []float64
		output []float64
	}{
		{[]float64{1.0, 1.0, 1.0, 1.0}, []float64{4.0, 0, 0, 0}},
	} {

		out := DCT1D(tt.input)
		pass := true

		if len(tt.output) != len(out) {
			t.Errorf("DCT1D(%v) is expected %v but got %v.", tt.input, tt.output, out)
		}

		for i := range out {
			if (out[i]-tt.output[i]) > EPSILON || (tt.output[i]-out[i]) > EPSILON {
				pass = false
			}
		}

		if !pass || len(tt.output) != len(out) {
			t.Errorf("DCT1D(%v) is expected %v but got %v.", tt.input, tt.output, out)
		}
	}
}

func TestDCT2D(t *testing.T) {
	for _, tt := range []struct {
		input  [][]float64
		output [][]float64
		w      int
		h      int
	}{
		{[][]float64{{1.0, 2.0, 3.0, 4.0},
			{5.0, 6.0, 7.0, 8.0},
			{9.0, 10.0, 11.0, 12.0},
			{13.0, 14.0, 15.0, 16.0}},
			[][]float64{{136.0, -12.6172881195958, 0.0, -0.8966830583359305},
				{-50.4691524783832, 0.0, 0.0, 0.0},
				{0.0, 0.0, 0.0, 0.0},
				{-3.586732233343722, 0.0, 0.0, 0.0}},
			4, 4},
	} {
		out := DCT2D(tt.input, tt.w, tt.h)
		pass := true

		for i := 0; i < tt.h; i++ {
			for j := 0; j < tt.w; j++ {
				if (out[i][j]-tt.output[i][j]) > EPSILON || (tt.output[i][j]-out[i][j]) > EPSILON {
					pass = false
				}
			}
		}

		if !pass {
			t.Errorf("DCT2D(%v) is expected %v but got %v.", tt.input, tt.output, out)
		}
	}
}

func BenchmarkForward(b *testing.B) {

	arr1 := make([]float64, 64)
	temp := make([]float64, 64)
	for i := 0; i < len(arr1); i++ {
		arr1[i] = float64(rand.Int63n(100000)) * 0.001
	}
	b.Run("forward", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			forwardTransform(arr1, temp, len(arr1))

		}
	})
	b.Run("forwardFast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			forwardTransformFast(arr1, temp, len(arr1))
		}
	})

	b.Run("forwardStatic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			forwardTransformStatic64(arr1)
		}
	})
	//

}
