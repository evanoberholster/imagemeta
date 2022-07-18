package transforms

func DCT1DFast64(input []float64) {
	var temp [64]float64
	for i := 0; i < 32; i++ {
		x, y := input[i], input[64-1-i]
		temp[i] = x + y
		temp[i+32] = (x - y) / dct64[i]
	}
	forwardTransformStatic32(temp[:32])
	forwardTransformStatic32(temp[32:])
	for i := 0; i < 32-1; i++ {
		input[i*2+0] = temp[i]
		input[i*2+1] = temp[i+32] + temp[i+32+1]
	}
	input[62], input[63] = temp[31], temp[63]
}

func forwardTransformStatic32(input []float64) {
	var temp [32]float64
	for i := 0; i < 16; i++ {
		x, y := input[i], input[32-1-i]
		temp[i] = x + y
		temp[i+16] = (x - y) / dct32[i]
	}
	forwardTransformStatic16(temp[:16])
	forwardTransformStatic16(temp[16:])
	for i := 0; i < 16-1; i++ {
		input[i*2+0] = temp[i]
		input[i*2+1] = temp[i+16] + temp[i+16+1]
	}

	input[30], input[31] = temp[15], temp[31]
}

func forwardTransformStatic16(input []float64) {
	var temp [16]float64
	for i := 0; i < 8; i++ {
		x, y := input[i], input[15-i]
		temp[i] = x + y
		temp[i+8] = (x - y) / dct16[i]
	}
	forwardDCT8(temp[:8])
	forwardDCT8(temp[8:])
	for i := 0; i < 8-1; i++ {
		input[i*2+0] = temp[i]
		input[i*2+1] = temp[i+8] + temp[i+8+1]
	}

	input[16-2], input[16-1] = temp[8-1], temp[16-1]
}

func forwardDCT8(input []float64) {
	x0 := input[0]
	x1 := input[1]
	x2 := input[2]
	x3 := input[3]
	x4 := input[4]
	x5 := input[5]
	x6 := input[6]
	x7 := input[7]

	tmp0 := x0 + x7
	tmp1 := x1 + x6
	tmp2 := x2 + x5
	tmp3 := x3 + x4

	tmp4 := (x0 - x7) / 1.9615705608064609
	tmp5 := (x1 - x6) / 1.6629392246050907
	tmp6 := (x2 - x5) / 1.1111404660392046
	tmp7 := (x3 - x4) / 0.3901806440322566

	a0 := tmp0 + tmp3
	a1 := tmp1 + tmp2
	a2 := (tmp0 - tmp3) / 1.8477590650225735
	a3 := (tmp1 - tmp2) / 0.7653668647301797

	a4 := tmp4 + tmp7
	a5 := tmp5 + tmp6
	a6 := (tmp4 - tmp7) / 1.8477590650225735
	a7 := (tmp5 - tmp6) / 0.7653668647301797

	b0 := (a0 - a1) / 1.4142135623730951
	b1 := (a4 - a5) / 1.4142135623730951
	b2 := (a2 - a3) / 1.4142135623730951
	b3 := (a6 - a7) / 1.4142135623730951

	input[0] = a0 + a1
	input[1] = a4 + a5 + a6 + a7 + b3
	input[2] = a2 + a3 + b2
	input[3] = a6 + a7 + b1 + b3
	input[4] = b0
	input[5] = b1 + b3
	input[6] = b2
	input[7] = b3
}
