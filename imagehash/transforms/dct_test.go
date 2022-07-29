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

func TestFastDCT2D(t *testing.T) {
	size := 64
	arr := make([]float64, size*size)
	arr2 := make([][]float64, size)

	for x := 0; x < size; x++ {
		arr2[x] = make([]float64, size)
		for i := 0; i < size; i++ {
			val := rand.Float64()
			arr[64*i+x] = val
			arr2[x][i] = val
		}
	}

	arr2 = DCT2D(arr2, 64, 64)
	dct2dFast(&arr)

	for x := 0; x < size; x++ {
		for i := 0; i < size; i++ {
			if float32(arr[64*i+x]) != float32(arr2[x][i]) {
				t.Error(arr[64*i+x], "!=", arr2[x][i])
			}
		}
	}
}

func TestForwardDC256(t *testing.T) {
	size := 256
	arr := make([]float64, size)
	arr2 := make([]float64, size)

	for x := 0; x < size; x++ {
		val := rand.Float64()
		arr[x] = val
		arr2[x] = val
	}
	temp := make([]float64, size)
	forwardTransform(arr, temp, len(arr))
	forwardDCT256(arr2)

	for i := 0; i < size; i++ {
		if arr[i] != arr2[i] {
			t.Error(i, arr[i], "!=", arr2[i])
		}
	}

}

func BenchmarkForwardDCT64(b *testing.B) {
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

	b.Run("forwardStatic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			forwardDCT64(arr1)
		}
	})
}

func dct2dFast(input *[]float64) {
	if len(*input) != 4096 {
		panic("Incorrect forward transform size")
	}
	for i := 0; i < 64; i++ { // height
		forwardDCT64((*input)[i*64 : 64*i+64])
	}

	var row [64]float64
	for i := 0; i < 64; i++ { // width
		for j := 0; j < 64; j++ {
			row[j] = (*input)[64*j+i]
		}
		forwardDCT64(row[:])
		for j := 0; j < 64; j++ {
			(*input)[64*j+i] = row[j]
		}
	}
}
