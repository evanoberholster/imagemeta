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

// asmDCT16 r1-r4, and 5 temporary registers
func asmDCT16(r1, r2, r3, r4 reg.Virtual) {
	tmp := make([]reg.VecVirtual, 4)
	for i := 0; i < len(tmp); i++ {
		tmp[i] = XMM()
	}

	x := XMM()

	Comment("DCT16")
	PSHUFD(U8(27), r3, x)
	VADDPS(x, r2, tmp[1])
	VSUBPS(x, r2, tmp[3])
	DIVPS(dct16values.Offset(4*4), tmp[3])

	PSHUFD(U8(27), r4, x)
	VADDPS(x, r1, tmp[0])
	VSUBPS(x, r1, tmp[2])
	DIVPS(dct16values.Offset(0), tmp[2])

	asmDCT8(tmp[0], tmp[1])
	asmDCT8(tmp[2], tmp[3])

	VPSRLDQ(U8(4), tmp[2], r1)
	VPSLLDQ(U8(12), tmp[3], r2)
	VADDPS(r1, r2, r1)
	VADDPS(r1, tmp[2], tmp[2])

	VPSRLDQ(U8(4), tmp[3], r2)
	VADDPS(r2, tmp[3], tmp[3])

	VUNPCKLPS(tmp[2], tmp[0], r1)
	VUNPCKHPS(tmp[2], tmp[0], r2)
	VUNPCKLPS(tmp[3], tmp[1], r3)
	VUNPCKHPS(tmp[3], tmp[1], r4)
	Comment("end DCT16")
}

// asmDCT32Unaligned
func asmDCT32Unaligned(mOffset []Mem) {
	Comment("DCT32 Unaligned")

	x := XMM()
	y := XMM()
	xmm := make([]reg.VecVirtual, 8)
	for i := 0; i < 4; i++ {
		xmm[i] = XMM()
		xmm[i+4] = XMM()
		MOVUPS(mOffset[i], x)
		MOVUPS(mOffset[7-i], y)
		PSHUFD(U8(27), y, y)
		VADDPS(y, x, xmm[i])
		VSUBPS(y, x, xmm[i+4])
		DIVPS(dct32values.Offset(i*4*4), xmm[i+4])
	}

	// DCT16
	asmDCT16(xmm[0], xmm[1], xmm[2], xmm[3])
	asmDCT16(xmm[4], xmm[5], xmm[6], xmm[7])
	//

	x1 := XMM()
	x2 := XMM()

	VPSRLDQ(U8(4), xmm[4], x1)
	VPSLLDQ(U8(12), xmm[5], x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, xmm[4], xmm[4])

	VUNPCKLPS(xmm[4], xmm[0], x1)
	MOVUPS(x1, mOffset[0])
	VUNPCKHPS(xmm[4], xmm[0], x2)
	MOVUPS(x2, mOffset[1])

	VPSRLDQ(U8(4), xmm[5], x1)
	VPSLLDQ(U8(12), xmm[6], x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, xmm[5], xmm[5])

	VUNPCKLPS(xmm[5], xmm[1], x1)
	MOVUPS(x1, mOffset[2])
	VUNPCKHPS(xmm[5], xmm[1], x2)
	MOVUPS(x2, mOffset[3])

	VPSRLDQ(U8(4), xmm[6], x1)
	VPSLLDQ(U8(12), xmm[7], x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, xmm[6], xmm[6])

	VUNPCKLPS(xmm[6], xmm[2], x1)
	MOVUPS(x1, mOffset[4])
	VUNPCKHPS(xmm[6], xmm[2], x2)
	MOVUPS(x2, mOffset[5])

	VPSRLDQ(U8(4), xmm[7], x2)
	VADDPS(x2, xmm[7], xmm[7])

	VUNPCKLPS(xmm[7], xmm[3], x1)
	MOVUPS(x1, mOffset[6])
	VUNPCKHPS(xmm[7], xmm[3], x2)
	MOVUPS(x2, mOffset[7])

	Comment("end DCT32 Unaligned")
}

// asmDCT32
func asmDCT32(mOffset []Mem) {

	Comment("DCT32")

	x := XMM()
	y := XMM()
	xmm := make([]reg.VecVirtual, 8)
	for i := 0; i < 4; i++ {
		xmm[i] = XMM()
		xmm[i+4] = XMM()
		MOVAPS(mOffset[i], x)
		PSHUFD(U8(27), mOffset[7-i], y)
		VADDPS(y, x, xmm[i])
		VSUBPS(y, x, xmm[i+4])
		DIVPS(dct32values.Offset(i*4*4), xmm[i+4])
	}

	// DCT16
	asmDCT16(xmm[0], xmm[1], xmm[2], xmm[3])
	asmDCT16(xmm[4], xmm[5], xmm[6], xmm[7])
	//

	x1 := XMM()
	x2 := XMM()

	VPSRLDQ(U8(4), xmm[4], x1)
	VPSLLDQ(U8(12), xmm[5], x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, xmm[4], xmm[4])

	VUNPCKLPS(xmm[4], xmm[0], x1)
	MOVAPS(x1, mOffset[0])
	VUNPCKHPS(xmm[4], xmm[0], x2)
	MOVAPS(x2, mOffset[1])

	VPSRLDQ(U8(4), xmm[5], x1)
	VPSLLDQ(U8(12), xmm[6], x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, xmm[5], xmm[5])

	VUNPCKLPS(xmm[5], xmm[1], x1)
	MOVAPS(x1, mOffset[2])
	VUNPCKHPS(xmm[5], xmm[1], x2)
	MOVAPS(x2, mOffset[3])

	VPSRLDQ(U8(4), xmm[6], x1)
	VPSLLDQ(U8(12), xmm[7], x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, xmm[6], xmm[6])

	VUNPCKLPS(xmm[6], xmm[2], x1)
	MOVAPS(x1, mOffset[4])
	VUNPCKHPS(xmm[6], xmm[2], x2)
	MOVAPS(x2, mOffset[5])

	VPSRLDQ(U8(4), xmm[7], x2)
	VADDPS(x2, xmm[7], xmm[7])

	VUNPCKLPS(xmm[7], xmm[3], x1)
	MOVAPS(x1, mOffset[6])
	VUNPCKHPS(xmm[7], xmm[3], x2)
	MOVAPS(x2, mOffset[7])

	Comment("end DCT32")
}

func asmDCT64(mOffset []Mem) {

	VZEROUPPER()

	ymm := make([]reg.VecVirtual, 8)
	for i := 0; i < len(ymm); i++ {
		ymm[i] = YMM()
	}

	ymmA := YMM()
	ymmB := YMM()
	ymmC := YMM()
	ymmD := YMM()
	perm := YMM()
	VPMOVZXBD(permValues.Offset(0), perm)

	VMOVAPS(mOffset[0], ymm[0])
	VMOVAPS(mOffset[2], ymm[1])
	VMOVAPS(mOffset[4], ymm[2])
	VMOVAPS(mOffset[6], ymm[3])
	VPERMD(mOffset[8], perm, ymm[4])
	VPERMD(mOffset[10], perm, ymm[5])
	VPERMD(mOffset[12], perm, ymm[6])
	VPERMD(mOffset[14], perm, ymm[7])
	VADDPS(ymm[7], ymm[0], ymmA)
	VADDPS(ymm[6], ymm[1], ymmB)
	VADDPS(ymm[5], ymm[2], ymmC)
	VADDPS(ymm[4], ymm[3], ymmD)
	VMOVAPS(ymmA, mOffset[0])
	VMOVAPS(ymmB, mOffset[2])
	VMOVAPS(ymmC, mOffset[4])
	VMOVAPS(ymmD, mOffset[6])

	VSUBPS(ymm[7], ymm[0], ymmA)
	VDIVPS(dct64values.Offset(0), ymmA, ymmA)
	VMOVAPS(ymmA, mOffset[8])
	VSUBPS(ymm[6], ymm[1], ymmB)
	VDIVPS(dct64values.Offset(4*8), ymmB, ymmB)
	VMOVAPS(ymmB, mOffset[10])
	VSUBPS(ymm[5], ymm[2], ymmC)
	VDIVPS(dct64values.Offset(4*16), ymmC, ymmC)
	VMOVAPS(ymmC, mOffset[12])
	VSUBPS(ymm[4], ymm[3], ymmD)
	VDIVPS(dct64values.Offset(4*24), ymmD, ymmD)
	VMOVAPS(ymmD, mOffset[14])

	VZEROUPPER()

	// DCT32
	asmDCT32([]Mem{mOffset[0], mOffset[1], mOffset[2], mOffset[3], mOffset[4], mOffset[5], mOffset[6], mOffset[7]})
	asmDCT32([]Mem{mOffset[8], mOffset[9], mOffset[10], mOffset[11], mOffset[12], mOffset[13], mOffset[14], mOffset[15]})
	//

	x1 := XMM()
	x2 := XMM()
	x3 := XMM()
	x4 := XMM()
	a1 := XMM()
	a2 := XMM()
	a3 := XMM()
	a4 := XMM()
	a5 := XMM()
	a6 := XMM()
	a7 := XMM()
	a8 := XMM()
	b1 := XMM()
	b2 := XMM()
	b3 := XMM()
	b4 := XMM()
	b5 := XMM()
	b6 := XMM()
	b7 := XMM()
	b8 := XMM()

	// 1-2
	MOVAPS(mOffset[0], a1)
	MOVAPS(mOffset[1], a2)
	MOVAPS(mOffset[2], a3)
	MOVAPS(mOffset[3], a4)

	MOVAPS(mOffset[8], b1)
	VPSRLDQ(U8(4), b1, x1)
	MOVAPS(mOffset[9], b2)
	VPSLLDQ(U8(12), b2, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b1, b1)

	VPSRLDQ(U8(4), b2, x1)
	MOVAPS(mOffset[10], b3)
	VPSLLDQ(U8(12), b3, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b2, b2)

	VUNPCKLPS(b1, a1, x1)
	MOVAPS(x1, mOffset[0])
	VUNPCKHPS(b1, a1, x2)
	MOVAPS(x2, mOffset[1])
	VUNPCKLPS(b2, a2, x3)
	MOVAPS(x3, mOffset[2])
	VUNPCKHPS(b2, a2, x4)
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

	MOVAPS(mOffset[4], a5)
	MOVAPS(mOffset[5], a6)
	MOVAPS(mOffset[6], a7)
	MOVAPS(mOffset[7], a8)

	VUNPCKLPS(b3, a3, x1)
	MOVAPS(x1, mOffset[4])
	VUNPCKHPS(b3, a3, x2)
	MOVAPS(x2, mOffset[5])
	VUNPCKLPS(b4, a4, x3)
	MOVAPS(x3, mOffset[6])
	VUNPCKHPS(b4, a4, x4)
	MOVAPS(x4, mOffset[7])

	// 5-6
	MOVAPS(mOffset[13], b6)
	MOVAPS(mOffset[14], b7)
	VPSRLDQ(U8(4), b5, x1)
	VPSLLDQ(U8(12), b6, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b5, b5)

	VUNPCKLPS(b5, a5, x1)
	MOVAPS(x1, mOffset[8])
	VUNPCKHPS(b5, a5, x2)
	MOVAPS(x2, mOffset[9])

	VPSRLDQ(U8(4), b6, x1)
	VPSLLDQ(U8(12), b7, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b6, b6)

	VUNPCKLPS(b6, a6, x3)
	MOVAPS(x3, mOffset[10])
	VUNPCKHPS(b6, a6, x4)
	MOVAPS(x4, mOffset[11])

	// 7-8

	MOVAPS(mOffset[15], b8)
	VPSRLDQ(U8(4), b7, x1)
	VPSLLDQ(U8(12), b8, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b7, b7)

	VUNPCKLPS(b7, a7, x1)
	MOVAPS(x1, mOffset[12])
	VUNPCKHPS(b7, a7, x2)
	MOVAPS(x2, mOffset[13])

	VPSRLDQ(U8(4), b8, x2)
	VADDPS(x2, b8, b8)

	VUNPCKLPS(b8, a8, x3)
	MOVAPS(x3, mOffset[14])
	VUNPCKHPS(b8, a8, x4)
	MOVAPS(x4, mOffset[15])
}

var (
	dct64values  = dct64()
	dct32values  = dct32()
	dct16values  = dct16()
	dct8values   = dct8()
	dct4values   = dct4()
	dct2values   = dct2()
	gatherValues = gather()
	permValues   = perm()
	xmm          = []reg.VecVirtual{}
)

func perm() Mem {
	val := GLOBL("perm", RODATA|NOPTR)
	DATA(0, U8(7))
	DATA(1, U8(6))
	DATA(2, U8(5))
	DATA(3, U8(4))
	DATA(4, U8(3))
	DATA(5, U8(2))
	DATA(6, U8(1))
	DATA(7, U8(0))
	return val
}

func gather() Mem {
	val := GLOBL("gather", RODATA|NOPTR)
	DATA(0, U32(64*0))
	DATA(4, U32(64*1))
	DATA(8, U32(64*2))
	DATA(12, U32(64*3))
	DATA(16, U32(64*4))
	DATA(20, U32(64*5))
	DATA(24, U32(64*6))
	DATA(28, U32(64*7))
	DATA(32, U32(1))
	return val
}

func asmDCT2D() {

	TEXT("asmDCT2D", NOSPLIT, "func(input []float32, tmp []float32) [64]float32")
	Pragma("noescape")
	Doc("asmDCT2D function returns a result of DCT2D by using the seperable property.\n  // DCT type II, unscaled. Algorithm by Byeong Gi Lee, 1984.\n // Custom built by Evan Oberholster for Hash64. Returns flattened pixels\n")

	input := Mem{Base: Load(Param("input").Base(), GP64())}
	ret0 := Mem{Base: Load(Param("tmp").Base(), GP64())}
	local := AllocLocal(64 * 4)
	ptr := AllocLocal(16)

	//ret0 := NewParamAddr("ret_0", 24) // UNSAFE method for returns

	// tmpOffset is input/output memory
	lOffset := make([]Mem, 16)
	for i := 0; i < len(lOffset); i++ {
		lOffset[i] = local.Offset(i * 4 * 4)
	}

	j := GP32()
	jdx := GP32()

	XORL(j, j) // set j to zero

	Label("j")
	CMPL(j, Imm(2))
	JE(LabelRef("continue"))
	Comment("Start innerloop instructions")
	// mOffset is input/output memory
	XORL(jdx, jdx) // set j to zero
	MOVL(U32(64), jdx)
	MULL(j)

	mOffset := make([]Mem, 16)
	for i := 0; i < len(mOffset); i++ {
		mOffset[i] = input.Offset(i*4*4).Idx(jdx, 4)
	}
	asmDCT64(mOffset)
	Comment("End innerloop instructions")

	INCL(j)
	JMP(LabelRef("j"))

	Label("continue")

	p := YMM()
	g := YMM()
	idx := YMM()
	mask := YMM()
	ymmA := YMM()
	ymmB := YMM()
	ymmC := YMM()
	ymmD := YMM()

	i := GP32()
	XORL(i, i) // set j to zero

	Label("i")
	CMPL(i, Imm(8))
	JE(LabelRef("done"))

	Comment("Start innerloop instructions")
	// Do something here
	Comment("--Loop load DCT64 values")
	VZEROUPPER()

	VPMOVZXBD(permValues.Offset(0), p)
	VMOVDQA(gatherValues.Offset(0), g)
	MOVL(i, ptr.Offset(0))
	VPBROADCASTD(ptr.Offset(0), idx)
	VPADDD(idx, g, g)

	VPCMPEQD(mask, mask, mask)
	VPGATHERDD(mask, input.Offset(0).Idx(g, 4), ymmA)
	VPCMPEQD(mask, mask, mask)
	VPGATHERDD(mask, input.Offset(64*8*4*7).Idx(g, 4), ymmB)
	VPERMD(ymmB, p, ymmB)
	VADDPS(ymmB, ymmA, ymmC)
	VMOVUPS(ymmC, lOffset[0])
	VSUBPS(ymmB, ymmA, ymmD)
	VDIVPS(dct64values.Offset(0), ymmD, ymmD)
	VMOVUPS(ymmD, lOffset[8])
	//
	VPCMPEQD(mask, mask, mask)
	VPGATHERDD(mask, input.Offset(64*8*4*1).Idx(g, 4), ymmA)
	VPCMPEQD(mask, mask, mask)
	VPGATHERDD(mask, input.Offset(64*8*4*6).Idx(g, 4), ymmB)
	VPERMD(ymmB, p, ymmB)
	VADDPS(ymmB, ymmA, ymmC)
	VMOVUPS(ymmC, lOffset[2])
	VSUBPS(ymmB, ymmA, ymmD)
	VDIVPS(dct64values.Offset(4*8), ymmD, ymmD)
	VMOVUPS(ymmD, lOffset[10])
	//
	VPCMPEQD(mask, mask, mask)
	VPGATHERDD(mask, input.Offset(64*8*4*2).Idx(g, 4), ymmA)
	VPCMPEQD(mask, mask, mask)
	VPGATHERDD(mask, input.Offset(64*8*4*5).Idx(g, 4), ymmB)
	VPERMD(ymmB, p, ymmB)
	VADDPS(ymmA, ymmB, ymmC)
	VMOVUPS(ymmC, lOffset[4])
	VSUBPS(ymmB, ymmA, ymmD)
	VDIVPS(dct64values.Offset(4*16), ymmD, ymmD)
	VMOVUPS(ymmD, lOffset[12])
	//
	VPCMPEQD(mask, mask, mask)
	VPGATHERDD(mask, input.Offset(64*8*4*3).Idx(g, 4), ymmA)
	VPCMPEQD(mask, mask, mask)
	VPGATHERDD(mask, input.Offset(64*8*4*4).Idx(g, 4), ymmB)
	VPERMD(ymmB, p, ymmB)
	VADDPS(ymmB, ymmA, ymmC)
	VMOVUPS(ymmC, lOffset[6])
	VSUBPS(ymmB, ymmA, ymmD)
	VDIVPS(dct64values.Offset(4*24), ymmD, ymmD)
	VMOVUPS(ymmD, lOffset[14])
	Comment("--Loop load DCT64 values")

	VZEROUPPER()

	asmDCT32Unaligned([]Mem{lOffset[0], lOffset[1], lOffset[2], lOffset[3], lOffset[4], lOffset[5], lOffset[6], lOffset[7]})
	asmDCT32Unaligned([]Mem{lOffset[8], lOffset[9], lOffset[10], lOffset[11], lOffset[12], lOffset[13], lOffset[14], lOffset[15]})

	B := XMM()
	C := XMM()
	// Forward DCT64 here

	MOVUPS(local.Offset(0), B)
	PEXTRD(Imm(0), B, ret0.Offset(0*8*4).Idx(i, 4))
	PEXTRD(Imm(1), B, ret0.Offset(2*8*4).Idx(i, 4))
	PEXTRD(Imm(2), B, ret0.Offset(4*8*4).Idx(i, 4))
	PEXTRD(Imm(3), B, ret0.Offset(6*8*4).Idx(i, 4))

	MOVUPS(local.Offset(32*4), B)
	MOVUPS(local.Offset(32*4+4), C)
	ADDPS(C, B)

	PEXTRD(Imm(0), B, ret0.Offset(1*8*4).Idx(i, 4))
	PEXTRD(Imm(1), B, ret0.Offset(3*8*4).Idx(i, 4))
	PEXTRD(Imm(2), B, ret0.Offset(5*8*4).Idx(i, 4))
	PEXTRD(Imm(3), B, ret0.Offset(7*8*4).Idx(i, 4))
	//
	Comment("End innerloop instructions")

	//ADDQ(Imm(4), j)
	INCL(i)
	JMP(LabelRef("i"))
	//

	Label("done")

	RET()

}

func asmForwardDCT64() {
	TEXT("asmForwardDCT64", NOSPLIT, "func(input []float32)")
	Doc("asmForwardDCT64 is a forward DCT transform for [64]float32")
	input := Mem{Base: Load(Param("input").Base(), GP64())}

	// mOffset is input/output memory
	mOffset := make([]Mem, 16)
	for i := 0; i < len(mOffset); i++ {
		mOffset[i] = input.Offset(i * 4 * 4)
	}

	asmDCT64(mOffset)

	RET()
}

func asmForwardDCT64Y() {
	TEXT("asmForwardDCT64", NOSPLIT, "func(input []float32)")
	//Pragma("noescape")
	Doc("asmForwardDCT64 is a forward DCT transform for [64]float32")
	input := Mem{Base: Load(Param("input").Base(), GP64())}

	VZEROUPPER()

	// mOffset is input/output memory
	mOffset := make([]Mem, 16)
	for i := 0; i < len(mOffset); i++ {
		mOffset[i] = input.Offset(i * 4 * 4)
	}

	x := XMM()
	y := XMM()
	b := XMM()

	for i := 0; i < 4; i++ {
		MOVAPS(mOffset[i], x)
		PSHUFD(U8(27), mOffset[15-i], y)
		VADDPS(x, y, b)
		MOVAPS(b, mOffset[i])
		VSUBPS(y, x, b)
		DIVPS(dct64values.Offset(i*4*4), b)
		PSHUFD(U8(27), mOffset[i+8], y)
		MOVAPS(b, mOffset[i+8])
		MOVAPS(mOffset[7-i], x)
		VADDPS(x, y, b)
		MOVAPS(b, mOffset[7-i])
		VSUBPS(y, x, b)
		DIVPS(dct64values.Offset((7-i)*4*4), b)
		MOVAPS(b, mOffset[15-i])
	}

	// DCT32
	asmDCT32([]Mem{mOffset[0], mOffset[1], mOffset[2], mOffset[3], mOffset[4], mOffset[5], mOffset[6], mOffset[7]})
	asmDCT32([]Mem{mOffset[8], mOffset[9], mOffset[10], mOffset[11], mOffset[12], mOffset[13], mOffset[14], mOffset[15]})
	//

	x1 := XMM()
	x2 := XMM()
	x3 := XMM()
	x4 := XMM()
	a1 := XMM()
	a2 := XMM()
	a3 := XMM()
	a4 := XMM()
	a5 := XMM()
	a6 := XMM()
	a7 := XMM()
	a8 := XMM()
	b1 := XMM()
	b2 := XMM()
	b3 := XMM()
	b4 := XMM()
	b5 := XMM()
	b6 := XMM()
	b7 := XMM()
	b8 := XMM()

	// 1-2
	MOVAPS(mOffset[0], a1)
	MOVAPS(mOffset[1], a2)
	MOVAPS(mOffset[2], a3)
	MOVAPS(mOffset[3], a4)

	MOVAPS(mOffset[8], b1)
	VPSRLDQ(U8(4), b1, x1)
	MOVAPS(mOffset[9], b2)
	VPSLLDQ(U8(12), b2, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b1, b1)

	VPSRLDQ(U8(4), b2, x1)
	MOVAPS(mOffset[10], b3)
	VPSLLDQ(U8(12), b3, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b2, b2)

	VUNPCKLPS(b1, a1, x1)
	MOVAPS(x1, mOffset[0])
	VUNPCKHPS(b1, a1, x2)
	MOVAPS(x2, mOffset[1])
	VUNPCKLPS(b2, a2, x3)
	MOVAPS(x3, mOffset[2])
	VUNPCKHPS(b2, a2, x4)
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

	MOVAPS(mOffset[4], a5)
	MOVAPS(mOffset[5], a6)
	MOVAPS(mOffset[6], a7)
	MOVAPS(mOffset[7], a8)

	VUNPCKLPS(b3, a3, x1)
	MOVAPS(x1, mOffset[4])
	VUNPCKHPS(b3, a3, x2)
	MOVAPS(x2, mOffset[5])
	VUNPCKLPS(b4, a4, x3)
	MOVAPS(x3, mOffset[6])
	VUNPCKHPS(b4, a4, x4)
	MOVAPS(x4, mOffset[7])

	// 5-6
	MOVAPS(mOffset[13], b6)
	MOVAPS(mOffset[14], b7)
	VPSRLDQ(U8(4), b5, x1)
	VPSLLDQ(U8(12), b6, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b5, b5)

	VUNPCKLPS(b5, a5, x1)
	MOVAPS(x1, mOffset[8])
	VUNPCKHPS(b5, a5, x2)
	MOVAPS(x2, mOffset[9])

	VPSRLDQ(U8(4), b6, x1)
	VPSLLDQ(U8(12), b7, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b6, b6)

	VUNPCKLPS(b6, a6, x3)
	MOVAPS(x3, mOffset[10])
	VUNPCKHPS(b6, a6, x4)
	MOVAPS(x4, mOffset[11])

	// 7-8

	MOVAPS(mOffset[15], b8)
	VPSRLDQ(U8(4), b7, x1)
	VPSLLDQ(U8(12), b8, x2)
	VADDPS(x1, x2, x1)
	VADDPS(x1, b7, b7)

	VUNPCKLPS(b7, a7, x1)
	MOVAPS(x1, mOffset[12])
	VUNPCKHPS(b7, a7, x2)
	MOVAPS(x2, mOffset[13])

	VPSRLDQ(U8(4), b8, x2)
	VADDPS(x2, b8, b8)

	VUNPCKLPS(b8, a8, x3)
	MOVAPS(x3, mOffset[14])
	VUNPCKHPS(b8, a8, x4)
	MOVAPS(x4, mOffset[15])

	RET()
}

func asmForwardDCT32() {
	TEXT("asmForwardDCT32", NOSPLIT, "func(input []float32)")
	Doc("asmForwardDCT32 is a forward DCT transform for [32]float32")
	input := Mem{Base: Load(Param("input").Base(), GP64())}

	mOffset := make([]Mem, 8)
	for i := 0; i < len(mOffset); i++ {
		mOffset[i] = input.Offset(i * 4 * 4)
	}

	asmDCT32(mOffset)

	RET()

	TEXT("asmForwardDCT16", NOSPLIT, "func(input []float32)")
	Doc("asmForwardDCT16 is a forward DCT transform for [16]float32")
	input = Mem{Base: Load(Param("input").Base(), GP64())}

	x1 := XMM()
	x2 := XMM()
	x3 := XMM()
	x4 := XMM()

	MOVAPS(input.Offset(0), x1)
	MOVAPS(input.Offset(4*4), x2)
	MOVAPS(input.Offset(4*8), x3)
	MOVAPS(input.Offset(4*12), x4)

	//asmDCT16(x1, x2, x3, x4)

	MOVAPS(x1, input.Offset(0))
	MOVAPS(x2, input.Offset(4*4))
	MOVAPS(x3, input.Offset(4*8))
	MOVAPS(x4, input.Offset(4*12))

	RET()
}

func main() {
	asmDCT2D()
	//asmForwardDCT32()
	asmForwardDCT64()

	//genPixelsToGray()
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
