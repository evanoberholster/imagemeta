package transforms

func DCT1DFast64(input []float64) {
	var temp [64]float64
	for i := 0; i < 32; i++ {
		x, y := input[i], input[63-i]
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
		x, y := input[i], input[31-i]
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
	forwardTransformStatic8(temp[:8])
	forwardTransformStatic8(temp[8:])
	for i := 0; i < 8-1; i++ {
		input[i*2+0] = temp[i]
		input[i*2+1] = temp[i+8] + temp[i+8+1]
	}

	input[16-2], input[16-1] = temp[8-1], temp[16-1]
}

func forwardTransformStatic8(input []float64) {
	var temp [8]float64
	x0, y0 := input[0], input[7]
	x1, y1 := input[1], input[6]
	x2, y2 := input[2], input[5]
	x3, y3 := input[3], input[4]

	temp[0] = x0 + y0
	temp[1] = x1 + y1
	temp[2] = x2 + y2
	temp[3] = x3 + y3
	temp[4] = (x0 - y0) / 1.9615705608064609
	temp[5] = (x1 - y1) / 1.6629392246050907
	temp[6] = (x2 - y2) / 1.1111404660392046
	temp[7] = (x3 - y3) / 0.3901806440322566

	forwardTransformStatic4(temp[:4])
	forwardTransformStatic4(temp[4:])

	input[0] = temp[0]
	input[1] = temp[4] + temp[5]
	input[2] = temp[1]
	input[3] = temp[5] + temp[6]
	input[4] = temp[2]
	input[5] = temp[6] + temp[7]
	input[6] = temp[3]
	input[7] = temp[7]
}

func forwardTransformStatic4(input []float64) {
	x0, x1, y1, y0 := input[0], input[1], input[2], input[3]

	t0 := x0 + y0
	t1 := x1 + y1
	t2 := (x0 - y0) / 1.8477590650225735
	t3 := (x1 - y1) / 0.7653668647301797

	input[0] = t0 + t1
	input[1] = t2 + t3 + (t2-t3)/1.4142135623730951
	input[2] = (t0 - t1) / 1.4142135623730951
	input[3] = (t2 - t3) / 1.4142135623730951
}
