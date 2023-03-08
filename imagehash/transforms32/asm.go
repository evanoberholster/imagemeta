//go:build ignore
// +build ignore

package main

import (
	"math"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
	//. "github.com/mmcloughlin/avo/reg"
)

//go:generate go run asm.go -out asm_x86.s -stubs stub.go

// Built with assistance from https://www.officedaytime.com/

func asmDCT8(r1, r2 reg.VecVirtual) {
	Comment("DCT8")
	x := XMM()
	y := XMM()
	a := XMM()
	b := XMM()

	PSHUFD(U8(27), r2, y)
	VADDPS(y, r1, a)
	VSUBPS(y, r1, b)
	DIVPS(dct8values, b)

	Comment("DCT4")
	PSHUFD(U8(180), a, a)
	PSHUFD(U8(180), b, b)

	VPUNPCKLDQ(b, a, x)
	VPUNPCKHDQ(b, a, y)
	VADDPS(y, x, a)
	VSUBPS(y, x, b)
	DIVPS(dct4values, b)

	VPUNPCKLDQ(b, a, x)
	VPUNPCKHDQ(b, a, y)

	Comment("DCT2")
	PSHUFD(U8(216), x, a)
	PSHUFD(U8(216), y, b)

	VADDPS(b, a, x)
	VSUBPS(b, a, b)
	DIVPS(dct2values, b)

	VADDPS(x, b, a)
	VBLENDPS(U8(3), x, a, a)

	VSHUFPS(U8(136), b, a, x)
	VSHUFPS(U8(221), b, a, y)

	VPSRLDQ(U8(4), y, b)
	ADDPS(b, y)
	VUNPCKLPS(y, x, r1)
	VUNPCKHPS(y, x, r2)
	Comment("end DCT8")
}

func asmDCT16(r1, r2, r3, r4 reg.Virtual) {
	Comment("DCT16")
	PSHUFD(U8(27), r3, y2)
	PSHUFD(U8(27), r4, y1)

	VADDPS(y1, r1, a1)
	VADDPS(y2, r2, a2)
	VSUBPS(y1, r1, b1)
	VSUBPS(y2, r2, b2)
	DIVPS(dct16values.Offset(0), b1)
	DIVPS(dct16values.Offset(4*4), b2)

	asmDCT8(a1, a2)
	asmDCT8(b1, b2)

	VPSRLDQ(U8(4), b1, r1)
	VPSLLDQ(U8(12), b2, r2)
	VADDPS(r1, r2, r1)
	VADDPS(r1, b1, b1)

	VPSRLDQ(U8(4), b2, r2)
	VADDPS(r2, b2, b2)

	VUNPCKLPS(b1, a1, r1)
	VUNPCKHPS(b1, a1, r2)
	VUNPCKLPS(b2, a2, r3)
	VUNPCKHPS(b2, a2, r4)
	Comment("end DCT16")
}

// asmDCT32o o1, o2, o3, o4, o5, o6, o7, o8
func asmDCT32(o1, o2, o3, o4, o5, o6, o7, o8 Mem) {
	r := make([]reg.VecVirtual, 8)
	for i := 0; i < len(r); i++ {
		r[i] = XMM()
	}

	Comment("DCT32")
	MOVAPS(o1, x1)
	MOVAPS(o2, x2)
	MOVAPS(o3, x3)
	MOVAPS(o4, x4)
	PSHUFD(U8(27), o5, y1)
	PSHUFD(U8(27), o6, y2)
	PSHUFD(U8(27), o7, y3)
	PSHUFD(U8(27), o8, y4)

	VADDPS(y4, x1, r[0])
	VADDPS(y3, x2, r[1])
	VADDPS(y2, x3, r[2])
	VADDPS(y1, x4, r[3])
	VSUBPS(y4, x1, r[4])
	VSUBPS(y3, x2, r[5])
	VSUBPS(y2, x3, r[6])
	VSUBPS(y1, x4, r[7])
	DIVPS(dct32values.Offset(0), r[4])
	DIVPS(dct32values.Offset(4*4), r[5])
	DIVPS(dct32values.Offset(4*8), r[6])
	DIVPS(dct32values.Offset(4*12), r[7])

	// DCT16
	asmDCT16(r[0], r[1], r[2], r[3])
	asmDCT16(r[4], r[5], r[6], r[7])
	//

	VPSRLDQ(U8(4), r[4], x1)
	VPSLLDQ(U8(12), r[5], x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, r[4], r[4])

	VPSRLDQ(U8(4), r[5], x1)
	VPSLLDQ(U8(12), r[6], x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, r[5], r[5])

	VPSRLDQ(U8(4), r[6], x1)
	VPSLLDQ(U8(12), r[7], x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, r[6], r[6])

	VPSRLDQ(U8(4), r[7], x2)
	VADDPS(x2, r[7], r[7])

	VUNPCKLPS(r[4], r[0], x1)
	VUNPCKHPS(r[4], r[0], x2)
	VUNPCKLPS(r[5], r[1], x3)
	VUNPCKHPS(r[5], r[1], x4)
	VUNPCKLPS(r[6], r[2], x5)
	VUNPCKHPS(r[6], r[2], x6)
	VUNPCKLPS(r[7], r[3], x7)
	VUNPCKHPS(r[7], r[3], x8)

	MOVAPS(x1, o1)
	MOVAPS(x2, o2)
	MOVAPS(x3, o3)
	MOVAPS(x4, o4)
	MOVAPS(x5, o5)
	MOVAPS(x6, o6)
	MOVAPS(x7, o7)
	MOVAPS(x8, o8)
	Comment("end DCT32")
}

var (
	dct64values = dct64()
	dct32values = dct32()
	dct16values = dct16()
	dct8values  = dct8()
	dct4values  = dct4()
	dct2values  = dct2()

	x1 = XMM()
	x2 = XMM()
	x3 = XMM()
	x4 = XMM()
	x5 = XMM()
	x6 = XMM()
	x7 = XMM()
	x8 = XMM()
	y1 = XMM()
	y2 = XMM()
	y3 = XMM()
	y4 = XMM()
	y5 = XMM()
	y6 = XMM()
	y7 = XMM()
	y8 = XMM()
	a1 = XMM()
	a2 = XMM()
	a3 = XMM()
	a4 = XMM()
	a5 = XMM()
	a6 = XMM()
	a7 = XMM()
	a8 = XMM()
	b1 = XMM()
	b2 = XMM()
	b3 = XMM()
	b4 = XMM()
	b5 = XMM()
	b6 = XMM()
	b7 = XMM()
	b8 = XMM()
)

func main() {

	TEXT("asmForwardDCT64", NOSPLIT, "func(input []float32)")
	//Pragma("noescape")
	Doc("asmForwardDCT64 is a forward DCT transform for [64]float32")
	m := Mem{Base: Load(Param("input").Base(), GP64())}

	VZEROUPPER()

	// mOffset is input/output memory
	mOffset := make([]Mem, 16)
	for i := 0; i < len(mOffset); i++ {
		mOffset[i] = m.Offset(i * 4 * 4)
	}

	MOVAPS(mOffset[0], x1)
	MOVAPS(mOffset[1], x2)
	MOVAPS(mOffset[2], x3)
	MOVAPS(mOffset[3], x4)
	PSHUFD(U8(27), mOffset[15], y1)
	PSHUFD(U8(27), mOffset[14], y2)
	PSHUFD(U8(27), mOffset[13], y3)
	PSHUFD(U8(27), mOffset[12], y4)
	VADDPS(x1, y1, a1)
	VADDPS(x2, y2, a2)
	VADDPS(x3, y3, a3)
	VADDPS(x4, y4, a4)
	MOVAPS(a1, mOffset[0])
	MOVAPS(a2, mOffset[1])
	MOVAPS(a3, mOffset[2])
	MOVAPS(a4, mOffset[3])
	VSUBPS(y1, x1, b1)
	VSUBPS(y2, x2, b2)
	VSUBPS(y3, x3, b3)
	VSUBPS(y4, x4, b4)
	DIVPS(dct64values.Offset(0), b1)
	DIVPS(dct64values.Offset(4*4), b2)
	DIVPS(dct64values.Offset(4*8), b3)
	DIVPS(dct64values.Offset(4*12), b4)

	MOVAPS(mOffset[4], x5)
	MOVAPS(mOffset[5], x6)
	MOVAPS(mOffset[6], x7)
	MOVAPS(mOffset[7], x8)
	PSHUFD(U8(27), mOffset[11], y5)
	PSHUFD(U8(27), mOffset[10], y6)
	PSHUFD(U8(27), mOffset[9], y7)
	PSHUFD(U8(27), mOffset[8], y8)
	VADDPS(x5, y5, a5)
	VADDPS(x6, y6, a6)
	VADDPS(x7, y7, a7)
	VADDPS(x8, y8, a8)
	MOVAPS(a5, mOffset[4])
	MOVAPS(a6, mOffset[5])
	MOVAPS(a7, mOffset[6])
	MOVAPS(a8, mOffset[7])
	MOVAPS(b1, mOffset[8])
	MOVAPS(b2, mOffset[9])
	MOVAPS(b3, mOffset[10])
	MOVAPS(b4, mOffset[11])
	VSUBPS(y5, x5, b5)
	VSUBPS(y6, x6, b6)
	VSUBPS(y7, x7, b7)
	VSUBPS(y8, x8, b8)
	DIVPS(dct64values.Offset(4*16), b5)
	DIVPS(dct64values.Offset(4*20), b6)
	DIVPS(dct64values.Offset(4*24), b7)
	DIVPS(dct64values.Offset(4*28), b8)
	MOVAPS(b5, mOffset[12])
	MOVAPS(b6, mOffset[13])
	MOVAPS(b7, mOffset[14])
	MOVAPS(b8, mOffset[15])

	// DCT32
	asmDCT32(mOffset[0], mOffset[1], mOffset[2], mOffset[3], mOffset[4], mOffset[5], mOffset[6], mOffset[7])
	asmDCT32(mOffset[8], mOffset[9], mOffset[10], mOffset[11], mOffset[12], mOffset[13], mOffset[14], mOffset[15])
	//

	// 1-2
	MOVAPS(mOffset[0], a1)
	MOVAPS(mOffset[1], a2)
	MOVAPS(mOffset[2], a3)
	MOVAPS(mOffset[3], a4)

	MOVAPS(mOffset[8], b1)
	MOVAPS(mOffset[9], b2)
	MOVAPS(mOffset[10], b3)
	VPSRLDQ(U8(4), b1, x1)
	VPSLLDQ(U8(12), b2, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b1, b1)

	VPSRLDQ(U8(4), b2, x1)
	VPSLLDQ(U8(12), b3, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b2, b2)

	VUNPCKLPS(b1, a1, x1)
	VUNPCKHPS(b1, a1, x2)
	VUNPCKLPS(b2, a2, x3)
	VUNPCKHPS(b2, a2, x4)
	MOVAPS(x1, mOffset[0])
	MOVAPS(x2, mOffset[1])
	MOVAPS(x3, mOffset[2])
	MOVAPS(x4, mOffset[3])

	// 3-4
	MOVAPS(mOffset[11], b4)
	MOVAPS(mOffset[12], b5)
	VPSRLDQ(U8(4), b3, x1)
	VPSLLDQ(U8(12), b4, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b3, b3)

	VPSRLDQ(U8(4), b4, x1)
	VPSLLDQ(U8(12), b5, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b4, b4)

	VUNPCKLPS(b3, a3, x5)
	VUNPCKHPS(b3, a3, x6)
	VUNPCKLPS(b4, a4, x7)
	VUNPCKHPS(b4, a4, x8)
	MOVAPS(mOffset[4], a5)
	MOVAPS(mOffset[5], a6)
	MOVAPS(mOffset[6], a7)
	MOVAPS(mOffset[7], a8)
	MOVAPS(x5, mOffset[4])
	MOVAPS(x6, mOffset[5])
	MOVAPS(x7, mOffset[6])
	MOVAPS(x8, mOffset[7])

	// 5-6
	MOVAPS(mOffset[13], b6)
	MOVAPS(mOffset[14], b7)
	VPSRLDQ(U8(4), b5, x1)
	VPSLLDQ(U8(12), b6, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b5, b5)

	VPSRLDQ(U8(4), b6, x1)
	VPSLLDQ(U8(12), b7, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b6, b6)

	VUNPCKLPS(b5, a5, y1)
	VUNPCKHPS(b5, a5, y2)
	VUNPCKLPS(b6, a6, y3)
	VUNPCKHPS(b6, a6, y4)
	MOVAPS(y1, mOffset[8])
	MOVAPS(y2, mOffset[9])
	MOVAPS(y3, mOffset[10])
	MOVAPS(y4, mOffset[11])

	// 7-8

	MOVAPS(mOffset[15], b8)
	VPSRLDQ(U8(4), b7, x1)
	VPSLLDQ(U8(12), b8, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b7, b7)

	VPSRLDQ(U8(4), b8, x2)
	VADDPS(x2, b8, b8)

	VUNPCKLPS(b7, a7, y5)
	VUNPCKHPS(b7, a7, y6)
	VUNPCKLPS(b8, a8, y7)
	VUNPCKHPS(b8, a8, y8)
	MOVAPS(y5, mOffset[12])
	MOVAPS(y6, mOffset[13])
	MOVAPS(y7, mOffset[14])
	MOVAPS(y8, mOffset[15])

	RET()

	TEXT("asmForwardDCT32", NOSPLIT, "func(input []float32)")
	Doc("asmForwardDCT32 is a forward DCT transform for [32]float32")
	m = Mem{Base: Load(Param("input").Base(), GP64())}

	mOffset = make([]Mem, 8)
	for i := 0; i < len(mOffset); i++ {
		mOffset[i] = m.Offset(i * 4 * 4)
	}
	asmDCT32(mOffset[0], mOffset[1], mOffset[2], mOffset[3], mOffset[4], mOffset[5], mOffset[6], mOffset[7])

	RET()

	TEXT("asmForwardDCT16", NOSPLIT, "func(input []float32)")
	Doc("asmForwardDCT16 is a forward DCT transform for [16]float32")
	m = Mem{Base: Load(Param("input").Base(), GP64())}

	MOVAPS(m.Offset(0), x1)
	MOVAPS(m.Offset(4*4), x2)
	MOVAPS(m.Offset(4*8), x3)
	MOVAPS(m.Offset(4*12), x4)

	asmDCT16(x1, x2, x3, x4)

	MOVAPS(x1, m.Offset(0))
	MOVAPS(x2, m.Offset(4*4))
	MOVAPS(x3, m.Offset(4*8))
	MOVAPS(x4, m.Offset(4*12))

	RET()

	genPixelsToGray()
	genPixelsToGray8()

	Generate()
}

func dct64() Mem {
	val := GLOBL("dct64", RODATA|NOPTR)
	for i := 0; i < 32; i++ {
		DATA(i*4, F32((math.Cos((float64(i)+0.5)*math.Pi/float64(64)) * 2)))
	}
	return val
}

func dct32() Mem {
	val := GLOBL("dct32", RODATA|NOPTR)
	for i := 0; i < 16; i++ {
		DATA(i*4, F32((math.Cos((float64(i)+0.5)*math.Pi/float64(32)) * 2)))
	}
	return val
}

func dct16() Mem {
	val := GLOBL("dct16", RODATA|NOPTR)
	for i := 0; i < 8; i++ {
		DATA(i*4, F32((math.Cos((float64(i)+0.5)*math.Pi/float64(16)) * 2)))
	}
	return val
}

func dct8() Mem {
	val := GLOBL("dct8", RODATA|NOPTR)
	for i := 0; i < 4; i++ {
		DATA(i*4, F32((math.Cos((float64(i)+0.5)*math.Pi/float64(8)) * 2)))
	}
	return val
}

func dct4() Mem {
	val := GLOBL("dct4", RODATA|NOPTR)
	DATA(0, F32((math.Cos((float64(0)+0.5)*math.Pi/float64(4)) * 2)))
	DATA(4, F32((math.Cos((float64(0)+0.5)*math.Pi/float64(4)) * 2)))
	DATA(8, F32((math.Cos((float64(1)+0.5)*math.Pi/float64(4)) * 2)))
	DATA(12, F32((math.Cos((float64(1)+0.5)*math.Pi/float64(4)) * 2)))
	return val
}

func dct2() Mem {
	val := GLOBL("dct2", RODATA|NOPTR)
	for i := 0; i < 4; i++ {
		DATA(i*4, F32((math.Cos((float64(0)+0.5)*math.Pi/float64(2)) * 2)))
	}
	return val
}

func xmm128() Mem {
	val := GLOBL("xmm128", RODATA|NOPTR)
	DATA(0, I32(128))
	DATA(4, I32(0x10101))
	DATA(8, I32(91881))
	DATA(12, I32(46802))
	DATA(16, I32(22554))
	DATA(20, I32(116130))
	DATA(24, F32(257.0))
	DATA(28, F32(float32(0.299*256/257))) // red
	DATA(32, F32(float32(0.587*256/257))) // green
	DATA(36, F32(float32(0.114)))         // blue
	return val
}

func ymm128() Mem {
	val := GLOBL("ymm128", RODATA|NOPTR)
	DATA(0, I32(128))
	DATA(4, I32(0x10101))
	DATA(8, I32(91881))
	DATA(12, I32(46802))
	DATA(16, I32(22554))
	DATA(20, I32(116130))
	DATA(24, F32(257.0))
	DATA(28, F32(float32(0.299*256/257))) // red
	DATA(32, F32(float32(0.587*256/257))) // green
	DATA(36, F32(float32(0.114)))         // blue
	return val
}

func genPixelsToGray() {
	constMem := xmm128()
	TEXT("AsmYCbCrToGray", NOSPLIT, "func(pixels []float32, minX, minY, maxX, maxY int, sY, sCb, sCr []uint8, yStride, cStride int)  uint64")
	Doc("AsmYCbCrToGray is a forward DCT transform for []float32")

	// Load 128 x4
	static := make([]reg.VecVirtual, 10)
	for i := 0; i < len(static); i++ {
		static[i] = XMM()
		VPBROADCASTD(constMem.Offset(i*4), static[i])
	}

	yStride := Load(Param("yStride"), GP64())
	idxyStrideBase := GP64()
	idxyStride := GP64()

	cStride := Load(Param("cStride"), GP64())
	idxcStrideBase := GP64()
	idxcStride := GP64()

	maxY := Load(Param("maxY"), GP64())
	maxX := Load(Param("maxX"), GP64())

	sY := Mem{Base: Load(Param("sY").Base(), GP64())}
	sCb := Mem{Base: Load(Param("sCb").Base(), GP64())}
	sCr := Mem{Base: Load(Param("sCr").Base(), GP64())}
	pixels := Load(Param("pixels").Base(), GP64())

	yyXMM := XMM()
	cbXMM := XMM()
	crXMM := XMM()

	rXMM := XMM()
	gXMM := XMM()
	bXMM := XMM()

	idx := GP64()
	XORQ(idx, idx)

	y := GP64()
	x := GP64()
	XORQ(y, y) // set y to zero
	XORQ(x, x) // set x to zero

	Label("y")
	CMPQ(y, maxY)
	JE(LabelRef("done"))

	// yStide Base
	MOVQ(yStride, idxyStrideBase)
	IMULQ(y, idxyStrideBase)

	// cStride Base
	MOVQ(cStride, idxcStrideBase)
	IMULQ(y, idxcStrideBase)

	Label("x")
	CMPQ(x, maxX)
	JE(LabelRef("xDone"))

	Comment("Start innerloop instructions")
	// Do something here

	// yStride
	MOVQ(idxyStrideBase, idxyStride)
	ADDQ(x, idxyStride)

	// cStride
	MOVQ(idxcStrideBase, idxcStride)
	ADDQ(x, idxcStride)

	VPMOVZXBD(sY.Idx(idxyStride, 1), yyXMM)
	VPMOVZXBD(sCb.Idx(idxcStride, 1), cbXMM)
	VPMOVZXBD(sCr.Idx(idxcStride, 1), crXMM)

	VPMULLD(static[1], yyXMM, yyXMM)
	VPSUBD(static[0], cbXMM, cbXMM)
	VPSUBD(static[0], crXMM, crXMM)

	// Red
	VPMULLD(static[2], crXMM, rXMM)
	VPADDD(yyXMM, rXMM, rXMM)
	VPSRAD(Imm(8), rXMM, rXMM)
	VCVTDQ2PS(rXMM, rXMM)
	//VDIVPS(static[6], gXMM, gXMM)
	VMULPS(static[7], rXMM, rXMM)

	// Green
	VPMULLD(static[3], crXMM, crXMM)
	VPMULLD(static[4], cbXMM, gXMM)
	VPSUBQ(gXMM, yyXMM, gXMM)
	VPSUBQ(crXMM, gXMM, gXMM)
	VPSRAD(Imm(8), gXMM, gXMM)
	VCVTDQ2PS(gXMM, gXMM)
	//VDIVPS(static[6], gXMM, gXMM)
	VMULPS(static[8], gXMM, gXMM)

	// Blue
	VPMULLD(static[5], cbXMM, bXMM)
	VPADDD(yyXMM, bXMM, bXMM)
	VPSRAD(Imm(8), bXMM, bXMM)
	VCVTDQ2PS(bXMM, bXMM)
	VMULPS(static[9], bXMM, bXMM)

	//r + b + g
	VADDPS(rXMM, bXMM, bXMM)
	VADDPS(gXMM, bXMM, bXMM)

	VMOVAPS(bXMM, Mem{Base: pixels, Index: idxyStride, Scale: 4})
	//
	Comment("End innerloop instructions")

	ADDQ(Imm(4), x)
	JMP(LabelRef("x"))

	Label("xDone")
	XORQ(x, x)
	INCQ(y)
	JMP(LabelRef("y"))

	Label("done")
	RET()
}

func genPixelsToGray8() {
	constMem := ymm128()
	TEXT("AsmYCbCrToGray8", NOSPLIT, "func(pixels []float32, minX, minY, maxX, maxY int, sY, sCb, sCr []uint8, yStride, cStride int)  uint64")
	Doc("AsmYCbCrToGray8 is a forward DCT transform for []float32")

	// Load 128 x4
	static := make([]reg.VecVirtual, 10)
	for i := 0; i < len(static); i++ {
		static[i] = YMM()
		VPBROADCASTD(constMem.Offset(i*4), static[i])
	}

	yStride := Load(Param("yStride"), GP64())
	idxyStrideBase := GP64()
	idxyStride := GP64()

	cStride := Load(Param("cStride"), GP64())
	idxcStrideBase := GP64()
	idxcStride := GP64()

	maxY := Load(Param("maxY"), GP64())
	maxX := Load(Param("maxX"), GP64())

	sY := Mem{Base: Load(Param("sY").Base(), GP64())}
	sCb := Mem{Base: Load(Param("sCb").Base(), GP64())}
	sCr := Mem{Base: Load(Param("sCr").Base(), GP64())}
	pixels := Load(Param("pixels").Base(), GP64())

	yyXMM := YMM()
	cbXMM := YMM()
	crXMM := YMM()

	rXMM := YMM()
	gXMM := YMM()
	bXMM := YMM()

	idx := GP64()
	XORQ(idx, idx)

	y := GP64()
	x := GP64()
	XORQ(y, y) // set y to zero
	XORQ(x, x) // set x to zero

	Label("y")
	CMPQ(y, maxY)
	JE(LabelRef("done"))

	// yStide Base
	MOVQ(yStride, idxyStrideBase)
	IMULQ(y, idxyStrideBase)

	// cStride Base
	MOVQ(cStride, idxcStrideBase)
	IMULQ(y, idxcStrideBase)

	Label("x")
	CMPQ(x, maxX)
	JE(LabelRef("xDone"))

	Comment("Start innerloop instructions")
	// Do something here

	// yStride
	MOVQ(idxyStrideBase, idxyStride)
	ADDQ(x, idxyStride)

	// cStride
	MOVQ(idxcStrideBase, idxcStride)
	ADDQ(x, idxcStride)

	VPMOVZXBD(sY.Idx(idxyStride, 1), yyXMM)
	VPMOVZXBD(sCb.Idx(idxcStride, 1), cbXMM)
	VPMOVZXBD(sCr.Idx(idxcStride, 1), crXMM)

	VPMULLD(static[1], yyXMM, yyXMM)
	VPSUBD(static[0], cbXMM, cbXMM)
	VPSUBD(static[0], crXMM, crXMM)

	// Red
	VPMULLD(static[2], crXMM, rXMM)
	VPADDD(yyXMM, rXMM, rXMM)
	VPSRAD(Imm(8), rXMM, rXMM)
	VCVTDQ2PS(rXMM, rXMM)
	//VDIVPS(static[6], gXMM, gXMM)
	VMULPS(static[7], rXMM, rXMM)

	// Green
	VPMULLD(static[3], crXMM, crXMM)
	VPMULLD(static[4], cbXMM, gXMM)
	VPSUBQ(gXMM, yyXMM, gXMM)
	VPSUBQ(crXMM, gXMM, gXMM)
	VPSRAD(Imm(8), gXMM, gXMM)
	VCVTDQ2PS(gXMM, gXMM)
	//VDIVPS(static[6], gXMM, gXMM)
	VMULPS(static[8], gXMM, gXMM)

	// Blue
	VPMULLD(static[5], cbXMM, bXMM)
	VPADDD(yyXMM, bXMM, bXMM)
	VPSRAD(Imm(8), bXMM, bXMM)
	VCVTDQ2PS(bXMM, bXMM)
	VMULPS(static[9], bXMM, bXMM)

	//r + b + g
	VADDPS(rXMM, bXMM, bXMM)
	VADDPS(gXMM, bXMM, bXMM)

	VMOVAPS(bXMM, Mem{Base: pixels, Index: idxyStride, Scale: 4})
	//
	Comment("End innerloop instructions")

	ADDQ(Imm(8), x)
	JMP(LabelRef("x"))

	Label("xDone")
	XORQ(x, x)
	INCQ(y)
	JMP(LabelRef("y"))

	Label("done")
	RET()
}
