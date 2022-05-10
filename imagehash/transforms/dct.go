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

func colDCT1D(wg *sync.WaitGroup, input *[]float64, i int) {
	temp := [pHashSize]float64{}
	forwardTransform((*input)[i*pHashSize:(i*pHashSize)+pHashSize], temp[:], len(temp))
	wg.Done()
}

func rowDCT1D(wg *sync.WaitGroup, input *[]float64, i int) {
	temp := [pHashSize]float64{}
	row := [pHashSize]float64{}
	for j := 0; j < pHashSize; j++ {
		row[j] = (*input)[i+(j*pHashSize)]
	}
	forwardTransform(row[:], temp[:], len(row[:]))
	for j := 0; j < len(row); j++ {
		(*input)[i+(j*pHashSize)] = row[j]
	}
	wg.Done()
}

// pHashSize is PHash Bitsize
const pHashSize = 64

// DCT2DFast function returns a result of DCT2D by using the seperable property.
func DCT2DFast(pixels *[]float64) {
	wg := new(sync.WaitGroup)
	for i := 0; i < pHashSize; i++ { // height
		wg.Add(1)
		//wp.DCT1DCol(wg, pixels, i)
		go colDCT1D(wg, pixels, i)
	}
	wg.Wait()

	for i := 0; i < pHashSize; i++ { // width
		wg.Add(1)
		//wp.DCT1DRow(wg, pixels, i)
		go rowDCT1D(wg, pixels, i)
	}
	wg.Wait()
}

func (wp *WorkerPool) DCT2DFast(pixels *[]float64) {
	wg := wp.wgPool.Get().(*sync.WaitGroup)
	defer wp.wgPool.Put(wg)
	for i := 0; i < pHashSize; i++ { // height
		wg.Add(1)
		wp.sendDCT1DCol(wg, pixels, i)
	}
	wg.Wait()

	for i := 0; i < pHashSize; i++ { // width
		wg.Add(1)
		wp.sendDCT1DRow(wg, pixels, i)
	}
	wg.Wait()
}