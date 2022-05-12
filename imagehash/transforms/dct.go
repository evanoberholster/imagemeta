// Package transforms provides the transformations for imagehash
package transforms

// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"math"
	"sync"
)

// DCT1D function returns result of DCT-II.
// DCT type II, unscaled. Algorithm by Byeong Gi Lee, 1984.
func DCT1D(input []float64) []float64 {
	temp := make([]float64, len(input))
	forwardTransform(input, temp, len(input))
	return input
}

func forwardTransform(input, temp []float64, Len int) {
	if Len == 1 {
		return
	}

	halfLen := Len / 2

	for i := 0; i < halfLen; i++ {
		x, y := input[i], input[Len-1-i]
		temp[i] = x + y
		temp[i+halfLen] = (x - y) / (math.Cos((float64(i)+0.5)*math.Pi/float64(Len)) * 2)
	}
	forwardTransform(temp, input, halfLen)
	forwardTransform(temp[halfLen:], input, halfLen)
	for i := 0; i < halfLen-1; i++ {
		input[i*2+0] = temp[i]
		input[i*2+1] = temp[i+halfLen] + temp[i+halfLen+1]
	}

	input[Len-2], input[Len-1] = temp[halfLen-1], temp[Len-1]
}

// DCT2D function returns a  result of DCT2D by using the seperable property.
func DCT2D(input [][]float64, w int, h int) [][]float64 {
	output := make([][]float64, h)
	for i := range output {
		output[i] = make([]float64, w)
	}

	wg := new(sync.WaitGroup)
	for i := 0; i < h; i++ {
		wg.Add(1)
		go func(i int) {
			output[i] = DCT1D(input[i])
			wg.Done()
		}(i)
	}

	wg.Wait()
	for i := 0; i < w; i++ {
		wg.Add(1)
		in := make([]float64, h)
		go func(i int) {
			for j := 0; j < h; j++ {
				in[j] = output[j][i]
			}
			rows := DCT1D(in)
			for j := 0; j < len(rows); j++ {
				output[j][i] = rows[j]
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	return output
}

// pHashSize is PHash Bitsize
const pHashSize = 64

// DCT2DFast function returns a result of DCT2D by using the seperable property.
// Fast version only works with pHashSize 64 will panic if another since is given.
func DCT2DFast(input *[]float64) {
	if len(*input) != 4096 {
		panic("Incorrect forward transform size")
	}
	for i := 0; i < pHashSize; i++ { // height
		DCT1DFast64((*input)[i*pHashSize : (i*pHashSize)+pHashSize])
	}

	for i := 0; i < pHashSize; i++ { // width
		row := [pHashSize]float64{}
		for j := 0; j < pHashSize; j++ {
			row[j] = (*input)[i+((j)*pHashSize)]
		}
		DCT1DFast64(row[:])
		for j := 0; j < len(row); j++ {
			(*input)[i+(j*pHashSize)] = row[j]
		}
	}
}

// DCT1DFast function returns result of DCT-II.
// DCT type II, unscaled. Algorithm by Byeong Gi Lee, 1984.
func DCT1DFast(input []float64) []float64 {
	temp := make([]float64, len(input))
	forwardTransform(input, temp, len(input))
	return input
}

func forwardTransformFast(input, temp []float64, Len int) {
	if Len == 1 {
		return
	}

	halfLen := Len / 2
	t := dctTables[halfLen>>1]
	for i := 0; i < halfLen; i++ {
		x, y := input[i], input[Len-1-i]
		temp[i] = x + y
		temp[i+halfLen] = (x - y) / t[i]
	}
	forwardTransformFast(temp, input, halfLen)
	forwardTransformFast(temp[halfLen:], input, halfLen)
	for i := 0; i < halfLen-1; i++ {
		input[i*2+0] = temp[i]
		input[i*2+1] = temp[i+halfLen] + temp[i+halfLen+1]
	}

	input[Len-2], input[Len-1] = temp[halfLen-1], temp[Len-1]
}

// DCT Tables
var (
	dctTables = [][]float64{
		dct2[:],  //0
		dct4[:],  //1
		dct8[:],  //2
		nil,      //3
		dct16[:], //4
		nil,      //5
		nil,      //6
		nil,      //7
		dct32[:], //8
		nil,      //9
		nil,      //10
		nil,      //11
		nil,      //12
		nil,      //13
		nil,      //14
		nil,      //15
		dct64[:], //16
	}
	dct64 = [32]float64{
		(math.Cos((float64(0)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(1)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(2)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(3)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(4)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(5)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(6)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(7)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(8)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(9)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(10)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(11)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(12)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(13)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(14)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(15)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(16)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(17)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(18)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(19)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(20)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(21)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(22)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(23)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(24)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(25)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(26)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(27)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(28)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(29)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(30)+0.5)*math.Pi/64) * 2),
		(math.Cos((float64(31)+0.5)*math.Pi/64) * 2),
	}
	dct32 = [16]float64{
		(math.Cos((float64(0)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(1)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(2)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(3)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(4)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(5)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(6)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(7)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(8)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(9)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(10)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(11)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(12)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(13)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(14)+0.5)*math.Pi/32) * 2),
		(math.Cos((float64(15)+0.5)*math.Pi/32) * 2),
	}
	dct16 = [8]float64{
		(math.Cos((float64(0)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(1)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(2)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(3)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(4)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(5)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(6)+0.5)*math.Pi/16) * 2),
		(math.Cos((float64(7)+0.5)*math.Pi/16) * 2),
	}
	dct8 = [4]float64{
		(math.Cos((float64(0)+0.5)*math.Pi/8) * 2),
		(math.Cos((float64(1)+0.5)*math.Pi/8) * 2),
		(math.Cos((float64(2)+0.5)*math.Pi/8) * 2),
		(math.Cos((float64(3)+0.5)*math.Pi/8) * 2),
	}
	dct4 = [2]float64{
		(math.Cos((float64(0)+0.5)*math.Pi/4) * 2),
		(math.Cos((float64(1)+0.5)*math.Pi/4) * 2),
	}
	dct2 = [1]float64{
		(math.Cos((float64(0)+0.5)*math.Pi/2) * 2),
	}
)
