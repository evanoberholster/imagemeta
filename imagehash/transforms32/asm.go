//go:build ignore
// +build ignore

package main

import (
	"math"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

//go:generate go run asm.go -out asm_x86.s -stubs stub.go

// Built with assistance from https://www.officedaytime.com/
func main() {
	asmDCT2D()
	asmForwardDCT64()
	asmForwardDCT256()

	yCbCrToGray()

	Generate()
}
func dct256() Mem {
	val := GLOBL("dct256", RODATA|NOPTR)
	for i := 0; i < 128; i++ {
		DATA(i*4, F32((math.Cos((float64(i)+0.5)*math.Pi/float64(256)) * 2)))
	}
	return val
}

func dct128() Mem {
	val := GLOBL("dct128", RODATA|NOPTR)
	for i := 0; i < 64; i++ {
		DATA(i*4, F32((math.Cos((float64(i)+0.5)*math.Pi/float64(128)) * 2)))
	}
	return val
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
	DATA(8, U8(1))
	DATA(9, U8(2))
	DATA(10, U8(3))
	DATA(11, U8(4))
	DATA(12, U8(5))
	DATA(13, U8(6))
	DATA(14, U8(7))
	DATA(15, U8(0))
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

func constPixelsToGrey8() Mem {
	val := GLOBL("constyCbCrGray", RODATA|NOPTR)
	DATA(0*4, I32(128))
	DATA(1*4, I32(0x10101))
	DATA(2*4, I32(91881))
	DATA(3*4, I32(46802))
	DATA(4*4, I32(22554))
	DATA(5*4, I32(116130))
	DATA(6*4, F32(float32(0.299*256/257))) // red
	DATA(7*4, F32(float32(0.587*256/257))) // green
	DATA(8*4, F32(float32(0.114)))         // blue
	return val
}

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

// asmDCT16 r1-r4 are XMM registers
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

// asmDCT32Unaligned optimized for 2nd DCT64 transformations
// Uses MOVUPS instead of MOVAPS due to local stack memory not being alligned.
// Only processes first 8 values.
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

	VPSRLDQ(U8(4), xmm[4], x)
	VPSLLDQ(U8(12), xmm[5], y)
	VADDPS(x, y, x)
	VADDPS(x, xmm[4], xmm[4])

	VUNPCKLPS(xmm[4], xmm[0], x)
	MOVUPS(x, mOffset[0])
	VUNPCKHPS(xmm[4], xmm[0], x)
	MOVUPS(x, mOffset[1])

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

	Comment("end DCT32")
}

func asmDCT64(mOffset []Mem) {
	VZEROUPPER()
	Comment("DCT64")
	ymm := make([]reg.VecVirtual, 8)
	for i := 0; i < len(ymm); i++ {
		ymm[i] = YMM()
	}

	ymmA := YMM()
	ymmB := YMM()
	ymmC := YMM()
	F := GP32()
	perm := YMM()
	VPMOVZXBD(permValues.Offset(0), perm)

	VMOVUPS(mOffset[0], ymm[0])
	VMOVUPS(mOffset[2], ymm[1])
	VMOVUPS(mOffset[4], ymm[2])
	VMOVUPS(mOffset[6], ymm[3])
	VPERMD(mOffset[8], perm, ymm[4])
	VPERMD(mOffset[10], perm, ymm[5])
	VPERMD(mOffset[12], perm, ymm[6])
	VPERMD(mOffset[14], perm, ymm[7])

	for i := 0; i < 4; i++ {
		VADDPS(ymm[7-i], ymm[i], ymmA)
		VMOVUPS(ymmA, mOffset[i*2])
		VSUBPS(ymm[7-i], ymm[i], ymmA)
		VDIVPS(dct64values.Offset(4*8*i), ymmA, ymmA)
		VMOVUPS(ymmA, mOffset[8+i*2])
	}

	VZEROALL()

	// DCT32
	asmDCT32(mOffset[:8])
	asmDCT32(mOffset[8:])
	//

	VZEROUPPER()

	MOVL(mOffset[15].Offset(3*4), F) // Copy last value to final memory

	VMOVUPS(mOffset[0], ymm[0])
	VMOVUPS(mOffset[2], ymm[1])
	VMOVUPS(mOffset[4], ymm[2])
	VMOVUPS(mOffset[6], ymm[3])

	for i := 0; i < 4; i++ {
		VMOVUPS(mOffset[8+i*2], ymmA)
		if i == 3 {
			VPMOVZXBD(permValues.Offset(8), perm)
			VPERMPS(mOffset[8+i*2], perm, ymmB)
			VADDPS(ymmB, ymmA, ymmB)
		} else {
			VADDPS(mOffset[8+i*2].Offset(4), ymmA, ymmB)
		}
		VUNPCKLPS(ymmB, ymm[i], ymmA)
		VUNPCKHPS(ymmB, ymm[i], ymmB)
		VPERM2F128(U8(2), ymmA, ymmB, ymmC)
		VMOVUPS(ymmC, mOffset[i*4])
		VPERM2F128(U8(19), ymmA, ymmB, ymmC)
		VMOVUPS(ymmC, mOffset[i*4+2])
	}

	MOVL(F, mOffset[15].Offset(3*4)) // Copy last value to final memory
	VZEROUPPER()
	Comment("end DCT64")
}

func asmDCT128(mOffset []Mem, lOffset []Mem) {
	Comment("DCT128")
	VZEROUPPER()

	ymm0 := YMM()
	ymm1 := YMM()
	ymmA := YMM()
	ymmB := YMM()
	ymmC := YMM()
	F := GP32()
	perm := YMM()
	VPMOVZXBD(permValues.Offset(0), perm)

	for i := 0; i < 8; i++ {
		VMOVUPS(mOffset[i*2], ymm0)
		VPERMD(mOffset[30-i*2], perm, ymm1)
		VADDPS(ymm1, ymm0, ymmA)
		VMOVUPS(ymmA, lOffset[i*2])
		VSUBPS(ymm1, ymm0, ymmA)
		VDIVPS(dct128values.Offset(4*8*i), ymmA, ymmA)
		VMOVUPS(ymmA, lOffset[16+i*2])
	}

	VZEROALL()

	// Perfom DCT64 twice on each half
	asmDCT64(lOffset[:16])
	asmDCT64(lOffset[16:])
	//

	VZEROUPPER()

	MOVL(lOffset[31].Offset(3*4), F) // Copy last value to final memory

	for i := 0; i < 8; i++ {
		VMOVUPS(lOffset[i*2], ymm0)
		VMOVUPS(lOffset[16+i*2], ymmA)
		if i == 7 { // if statement to avoid memory overrun
			VPMOVZXBD(permValues.Offset(8), perm)
			VPERMPS(lOffset[16+i*2], perm, ymmB)
			VADDPS(ymmB, ymmA, ymmB)
		} else {
			VADDPS(lOffset[16+i*2].Offset(4), ymmA, ymmB)
		}
		VUNPCKLPS(ymmB, ymm0, ymmA)
		VUNPCKHPS(ymmB, ymm0, ymmB)
		VPERM2F128(U8(2), ymmA, ymmB, ymmC)
		VMOVUPS(ymmC, mOffset[i*4])
		VPERM2F128(U8(19), ymmA, ymmB, ymmC)
		VMOVUPS(ymmC, mOffset[i*4+2])
	}

	MOVL(F, mOffset[31].Offset(3*4)) // Copy last value to final memory
	VZEROUPPER()
	Comment("end DCT128")
}

func asmDCT256(mOffset []Mem, lOffset []Mem) {
	Comment("DCT256")
	VZEROUPPER()

	ymm0 := YMM()
	ymm1 := YMM()
	ymmA := YMM()
	ymmB := YMM()
	ymmC := YMM()
	F := GP32()
	perm := YMM()
	VPMOVZXBD(permValues.Offset(0), perm)

	for i := 0; i < 16; i++ {
		VMOVUPS(mOffset[i*2], ymm0)
		VPERMD(mOffset[62-i*2], perm, ymm1)
		VADDPS(ymm1, ymm0, ymmA)
		VMOVUPS(ymmA, lOffset[i*2])
		VSUBPS(ymm1, ymm0, ymmA)
		VDIVPS(dct256values.Offset(4*8*i), ymmA, ymmA)
		VMOVUPS(ymmA, lOffset[32+i*2])
	}

	VZEROALL()

	// DCT32
	asmDCT128(lOffset[:32], mOffset[:32])
	asmDCT128(lOffset[32:], mOffset[32:])
	//

	VZEROUPPER()

	MOVL(lOffset[63].Offset(3*4), F) // Copy last value to final memory

	for i := 0; i < 16; i++ {
		VMOVUPS(lOffset[i*2], ymm0)
		VMOVUPS(lOffset[32+i*2], ymmA)
		if i == 15 { // if statement to avoid memory overrun
			VPMOVZXBD(permValues.Offset(8), perm)
			VPERMPS(lOffset[32+i*2], perm, ymmB)
			VADDPS(ymmB, ymmA, ymmB)
		} else {
			VADDPS(lOffset[32+i*2].Offset(4), ymmA, ymmB)
		}
		VUNPCKLPS(ymmB, ymm0, ymmA)
		VUNPCKHPS(ymmB, ymm0, ymmB)
		VPERM2F128(U8(2), ymmA, ymmB, ymmC)
		VMOVUPS(ymmC, mOffset[i*4])
		VPERM2F128(U8(19), ymmA, ymmB, ymmC)
		VMOVUPS(ymmC, mOffset[i*4+2])
	}

	MOVL(F, mOffset[63].Offset(3*4)) // Copy last value to final memory
	VZEROUPPER()
	Comment("end DCT256")
}

func asmDCT2D() {
	TEXT("asmDCT2DHash64", NOSPLIT|NOPTR, "func(input []float32) [64]float32")
	Pragma("noescape")
	Doc("asmDCT2DHash64 function returns a result of DCT2D by using the seperable property.\n  // DCT type II, unscaled. Algorithm by Byeong Gi Lee, 1984.\n // Custom built by Evan Oberholster for Hash64. Returns flattened pixels\n")

	input := Mem{Base: Load(Param("input").Base(), GP64())}
	local := AllocLocal(64 * 4)
	ptr := AllocLocal(16)

	// UNSAFE method for returns
	// ret0 is the first Returned parameter [64]float32.
	retBase := 24
	ret0 := NewParamAddr("ret_0", retBase)
	ret8 := NewParamAddr("ret_8", retBase+8*4*1)
	ret16 := NewParamAddr("ret_16", retBase+8*4*2)
	ret24 := NewParamAddr("ret_24", retBase+8*4*3)
	ret32 := NewParamAddr("ret_32", retBase+8*4*4)
	ret40 := NewParamAddr("ret_40", retBase+8*4*5)
	ret48 := NewParamAddr("ret_48", retBase+8*4*6)
	ret56 := NewParamAddr("ret_56", retBase+8*4*7)

	// lOffset is local stack allocated memory for efficient copying of tempory values
	lOffset := make([]Mem, 16)
	for i := 0; i < len(lOffset); i++ {
		lOffset[i] = local.Offset(i * 4 * 4)
	}

	i := GP32()
	XORL(i, i) // set j to zero
	j := GP32()
	XORL(j, j) // set j to zero
	jdx := GP32()
	XORL(jdx, jdx) // set jdx to zero

	mOffset := make([]Mem, 16)
	for i := 0; i < len(mOffset); i++ {
		mOffset[i] = input.Offset(i*4*4).Idx(jdx, 4)
	}

	A := XMM()
	B := XMM()
	p := YMM()
	g := YMM()
	idx := YMM()
	mask := YMM()
	ymmA := YMM()
	ymmB := YMM()
	ymmC := YMM()

	Label("j") // J loop for first DCT
	CMPL(j, Imm(64))
	JE(LabelRef("i"))

	Comment("Start innerloop instructions")
	MOVL(U32(64), jdx)
	IMULL(j, jdx) // Calculate j Index. jdx = 64 * j

	asmDCT64(mOffset) // Perform first DCT64 transformation
	Comment("End innerloop instructions")

	INCL(j)
	JMP(LabelRef("j")) // End of J loop

	Label("i") // I Loop for second DCT64
	CMPL(i, Imm(8))
	JE(LabelRef("done"))

	Comment("Start innerloop instructions")
	Comment("--Loop load DCT64 values")
	MOVL(i, ptr.Offset(0))                 // Move i to memory to vaoid using AVX512 instructions
	VZEROUPPER()                           // ZERO upper bits of YMM registers for performance
	VPMOVZXBD(permValues.Offset(0), p)     // Permutation values for reversal of YMM values
	VPBROADCASTD(ptr.Offset(0), idx)       // Register GP32 Broadcast to YMM requires V5+VL. To limit to ACX2 had to use stack memory.
	VPADDD(gatherValues.Offset(0), idx, g) // Add current index to baseline gather offsets

	for i := 0; i < 4; i++ {
		VPCMPEQD(mask, mask, mask)
		VPGATHERDD(mask, input.Offset(64*8*4*i).Idx(g, 4), ymmA)
		VPCMPEQD(mask, mask, mask)
		VPGATHERDD(mask, input.Offset(64*8*4*(7-i)).Idx(g, 4), ymmB)
		VPERMD(ymmB, p, ymmB)
		VADDPS(ymmB, ymmA, ymmC)
		VMOVUPS(ymmC, lOffset[i*2])
		VSUBPS(ymmB, ymmA, ymmC)
		VDIVPS(dct64values.Offset(i*4*8), ymmC, ymmC)
		VMOVUPS(ymmC, lOffset[i*2+8])
	}
	Comment("--Loop load DCT64 values")

	VZEROUPPER() // ZERO upper bits of YMM

	// Perfom DCT32 with Unaligned memory due to Local Stack Allocated memory
	asmDCT32Unaligned(lOffset[:8])
	asmDCT32Unaligned(lOffset[8:])

	MOVUPS(lOffset[0], A)             // Move from local to XMM
	PEXTRD(Imm(0), A, ret0.Idx(i, 4)) // Uses SSE 4.1 to extract each item from XMM
	PEXTRD(Imm(1), A, ret16.Idx(i, 4))
	PEXTRD(Imm(2), A, ret32.Idx(i, 4))
	PEXTRD(Imm(3), A, ret48.Idx(i, 4))

	MOVUPS(lOffset[8], A)           // Move from local to XMM
	MOVUPS(lOffset[8].Offset(4), B) // Move from local to XMM
	ADDPS(A, B)

	PEXTRD(Imm(0), B, ret8.Idx(i, 4))
	PEXTRD(Imm(1), B, ret24.Idx(i, 4))
	PEXTRD(Imm(2), B, ret40.Idx(i, 4))
	PEXTRD(Imm(3), B, ret56.Idx(i, 4))

	Comment("End innerloop instructions")

	INCL(i)
	JMP(LabelRef("i"))
	//

	Label("done")

	RET()

}
func asmForwardDCT256() {
	TEXT("asmForwardDCT256", NOPTR, "func(input []float32)")
	Doc("asmForwardDCT256 is a forward DCT transform for [256]float32")
	input := Mem{Base: Load(Param("input").Base(), GP64())}
	local := AllocLocal(4 * 256)
	// mOffset is input/output memory
	mOffset := make([]Mem, 64)
	for i := 0; i < len(mOffset); i++ {
		mOffset[i] = input.Offset(i * 4 * 4)
	}
	lOffset := make([]Mem, 64)
	for i := 0; i < len(lOffset); i++ {
		lOffset[i] = local.Offset(i * 4 * 4)
	}
	asmDCT256(mOffset, lOffset)

	RET()
}

func asmForwardDCT64() {
	TEXT("asmForwardDCT64", NOSPLIT|NOPTR, "func(input []float32)")
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

var (
	dct256values = dct256()
	dct128values = dct128()
	dct64values  = dct64()
	dct32values  = dct32()
	dct16values  = dct16()
	dct8values   = dct8()
	dct4values   = dct4()
	dct2values   = dct2()
	gatherValues = gather()
	permValues   = perm()

	// Pixels to Gray (YMM)
	pixelsToGray8Values = constPixelsToGrey8()
)

func yCbCrToGray() {
	TEXT("asmYCbCrToGray", NOSPLIT|NOPTR, "func(pixels []float32, minX, minY, maxX, maxY int, sY, sCb, sCr []uint8, yStride, cStride int)")
	Doc("asmYCbCrToGra converts a YCbCr image to grayscale pixels. \n // Converts using 8x SIMD instructions and requires AVX and AVX2 ")

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

	VZEROUPPER()
	static := make([]reg.VecVirtual, 9)
	for i := 0; i < len(static); i++ {
		static[i] = YMM()
		VPBROADCASTD(pixelsToGray8Values.Offset(i*4), static[i])
	}

	yy := YMM()
	Cb := YMM()
	Cr := YMM()

	red := YMM()
	green := YMM()
	blue := YMM()
	gray := YMM()

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

	// yStride
	MOVQ(idxyStrideBase, idxyStride)
	ADDQ(x, idxyStride)

	// cStride
	MOVQ(idxcStrideBase, idxcStride)
	ADDQ(x, idxcStride)

	VPMOVZXBD(sY.Idx(idxyStride, 1), yy)
	VPMOVZXBD(sCb.Idx(idxcStride, 1), Cb)
	VPMOVZXBD(sCr.Idx(idxcStride, 1), Cr)

	VPMULLD(static[1], yy, yy)
	VPSUBD(static[0], Cb, Cb)
	VPSUBD(static[0], Cr, Cr)

	// Red
	VPMULLD(static[2], Cr, red)
	VPADDD(yy, red, red)
	VPSRAD(Imm(8), red, red) // Divide by 256 (Shift 8 bytes right)
	VCVTDQ2PS(red, red)
	VMULPS(static[6], red, red) // Multiply by adjusting factor

	// Green
	VPMULLD(static[3], Cr, Cr)
	VPMULLD(static[4], Cb, green)
	VPSUBQ(green, yy, green)
	VPSUBQ(Cr, green, green)
	VPSRAD(Imm(8), green, green) // Divide by 256 (Shift 8 bytes right)
	VCVTDQ2PS(green, green)
	VMULPS(static[7], green, green) // Multiply by adjusting factor

	// Blue
	VPMULLD(static[5], Cb, blue)
	VPADDD(yy, blue, blue)
	VPSRAD(Imm(8), blue, blue) // Divide by 256 (Shift 8 bytes right)
	VCVTDQ2PS(blue, blue)
	VMULPS(static[8], blue, blue) // Multiply by adjusting factor

	// Add Red + Blue + Green
	VADDPS(red, blue, gray)
	VADDPS(green, gray, gray)

	// Move result to memory
	VMOVAPS(gray, Mem{Base: pixels, Index: idxyStride, Scale: 4})
	Comment("End innerloop instructions")

	ADDQ(Imm(8), x)
	JMP(LabelRef("x"))

	Label("xDone")
	XORQ(x, x)
	INCQ(y)
	JMP(LabelRef("y"))

	Label("done")
	VZEROUPPER()
	RET()
}
