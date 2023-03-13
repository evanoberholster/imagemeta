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

// asmDCT6W4ithIndex adjusts for index offsets with jds
func asmDCT64WithIndex(mOffset []Mem, mm Mem, jdx reg.GPVirtual) {
	VZEROUPPER()

	ymm := make([]reg.VecVirtual, 8)
	for i := 0; i < len(ymm); i++ {
		ymm[i] = YMM()
	}

	ymmA := YMM()
	ymmB := YMM()
	ymmC := YMM()
	perm := YMM()
	F := GP32()
	VPMOVZXBD(permValues.Offset(0), perm)

	VMOVAPS(mOffset[0], ymm[0])
	VMOVAPS(mOffset[2], ymm[1])
	VMOVAPS(mOffset[4], ymm[2])
	VMOVAPS(mOffset[6], ymm[3])
	VPERMD(mOffset[8], perm, ymm[4])
	VPERMD(mOffset[10], perm, ymm[5])
	VPERMD(mOffset[12], perm, ymm[6])
	VPERMD(mOffset[14], perm, ymm[7])

	for i := 0; i < 4; i++ {
		VADDPS(ymm[7-i], ymm[i], ymmA)
		VMOVAPS(ymmA, mOffset[i*2])
		VSUBPS(ymm[7-i], ymm[i], ymmA)
		VDIVPS(dct64values.Offset(4*8*i), ymmA, ymmA)
		VMOVAPS(ymmA, mOffset[8+i*2])
	}
	VZEROUPPER()

	// DCT32
	asmDCT32([]Mem{mOffset[0], mOffset[1], mOffset[2], mOffset[3], mOffset[4], mOffset[5], mOffset[6], mOffset[7]})
	asmDCT32([]Mem{mOffset[8], mOffset[9], mOffset[10], mOffset[11], mOffset[12], mOffset[13], mOffset[14], mOffset[15]})
	//

	VZEROUPPER()

	VMOVUPS(mOffset[0], ymm[0])
	VMOVUPS(mOffset[2], ymm[1])
	VMOVUPS(mOffset[4], ymm[2])
	VMOVUPS(mOffset[6], ymm[3])
	MOVL(mm.Offset(63*4).Idx(jdx, 4), F) // Copy last value to final memory

	for i := 0; i < 4; i++ {
		VMOVUPS(mOffset[8+i*2], ymmA)
		VADDPS(mm.Offset((8+i*2)*4*4+4).Idx(jdx, 4), ymmA, ymmB)
		VUNPCKLPS(ymmB, ymm[i], ymmA)
		VUNPCKHPS(ymmB, ymm[i], ymmB)
		VPERM2F128(U8(2), ymmA, ymmB, ymmC)
		VMOVUPS(ymmC, mOffset[i*4])
		VPERM2F128(U8(19), ymmA, ymmB, ymmC)
		VMOVUPS(ymmC, mOffset[i*4+2])
	}

	MOVL(F, mm.Offset(63*4).Idx(jdx, 4)) // Copy last value to final memory

	VZEROUPPER()
}

func asmDCT64(mOffset []Mem, mm Mem) {
	VZEROUPPER()
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

	VMOVAPS(mOffset[0], ymm[0])
	VMOVAPS(mOffset[2], ymm[1])
	VMOVAPS(mOffset[4], ymm[2])
	VMOVAPS(mOffset[6], ymm[3])
	VPERMD(mOffset[8], perm, ymm[4])
	VPERMD(mOffset[10], perm, ymm[5])
	VPERMD(mOffset[12], perm, ymm[6])
	VPERMD(mOffset[14], perm, ymm[7])

	for i := 0; i < 4; i++ {
		VADDPS(ymm[7-i], ymm[i], ymmA)
		VMOVAPS(ymmA, mOffset[i*2])
		VSUBPS(ymm[7-i], ymm[i], ymmA)
		VDIVPS(dct64values.Offset(4*8*i), ymmA, ymmA)
		VMOVAPS(ymmA, mOffset[8+i*2])
	}

	VZEROALL()

	// DCT32
	asmDCT32([]Mem{mOffset[0], mOffset[1], mOffset[2], mOffset[3], mOffset[4], mOffset[5], mOffset[6], mOffset[7]})
	asmDCT32([]Mem{mOffset[8], mOffset[9], mOffset[10], mOffset[11], mOffset[12], mOffset[13], mOffset[14], mOffset[15]})
	//

	VZEROUPPER()

	MOVL(mm.Offset(63*4), F) // Copy last value to final memory

	VMOVUPS(mOffset[0], ymm[0])
	VMOVUPS(mOffset[2], ymm[1])
	VMOVUPS(mOffset[4], ymm[2])
	VMOVUPS(mOffset[6], ymm[3])

	for i := 0; i < 4; i++ {
		VMOVUPS(mOffset[8+i*2], ymmA)
		VADDPS(mm.Offset((8+i*2)*4*4+4), ymmA, ymmB)
		VUNPCKLPS(ymmB, ymm[i], ymmA)
		VUNPCKHPS(ymmB, ymm[i], ymmB)
		VPERM2F128(U8(2), ymmA, ymmB, ymmC)
		VMOVUPS(ymmC, mOffset[i*4])
		VPERM2F128(U8(19), ymmA, ymmB, ymmC)
		VMOVUPS(ymmC, mOffset[i*4+2])
	}

	MOVL(F, mm.Offset(63*4)) // Copy last value to final memory
	VZEROUPPER()
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

	// Pixels to Gray (XMM)
	pixelsToGray4Values = constPixelsToGrey4()

	// Pixels to Gray (YMM)
	pixelsToGray8Values = constPixelsToGrey8()
)

func asmDCT2D() {
	TEXT("asmDCT2DHash64", NOSPLIT, "func(input []float32) [64]float32")
	Pragma("noescape")
	Doc("asmDCT2DHash64 function returns a result of DCT2D by using the seperable property.\n  // DCT type II, unscaled. Algorithm by Byeong Gi Lee, 1984.\n // Custom built by Evan Oberholster for Hash64. Returns flattened pixels\n")

	input := Mem{Base: Load(Param("input").Base(), GP64())}
	local := AllocLocal(64 * 4)
	ptr := AllocLocal(16)

	// ret0 is the first Returned parameter [64]float32.
	ret0 := NewParamAddr("ret_0", 24) // UNSAFE method for returns

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

	asmDCT64WithIndex(mOffset, input, jdx) // Perform first DCT64 transformation
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
	asmDCT32Unaligned([]Mem{lOffset[0], lOffset[1], lOffset[2], lOffset[3], lOffset[4], lOffset[5], lOffset[6], lOffset[7]})
	asmDCT32Unaligned([]Mem{lOffset[8], lOffset[9], lOffset[10], lOffset[11], lOffset[12], lOffset[13], lOffset[14], lOffset[15]})

	MOVUPS(lOffset[0], A)                           // Move from local to XMM
	PEXTRD(Imm(0), A, ret0.Offset(0*8*4).Idx(i, 4)) // Uses SSE 4.1 to extract each item from XMM
	PEXTRD(Imm(1), A, ret0.Offset(2*8*4).Idx(i, 4))
	PEXTRD(Imm(2), A, ret0.Offset(4*8*4).Idx(i, 4))
	PEXTRD(Imm(3), A, ret0.Offset(6*8*4).Idx(i, 4))

	MOVUPS(local.Offset(32*4), A)   // Move from local to XMM
	MOVUPS(local.Offset(32*4+4), B) // Move from local to XMM
	ADDPS(A, B)

	PEXTRD(Imm(0), B, ret0.Offset(1*8*4).Idx(i, 4))
	PEXTRD(Imm(1), B, ret0.Offset(3*8*4).Idx(i, 4))
	PEXTRD(Imm(2), B, ret0.Offset(5*8*4).Idx(i, 4))
	PEXTRD(Imm(3), B, ret0.Offset(7*8*4).Idx(i, 4))

	Comment("End innerloop instructions")

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

	asmDCT64(mOffset, input)

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

func constPixelsToGrey4() Mem {
	val := GLOBL("pg4", RODATA|NOPTR)
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

func constPixelsToGrey8() Mem {
	val := GLOBL("pg8", RODATA|NOPTR)
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
	TEXT("AsmYCbCrToGray", NOSPLIT, "func(pixels []float32, minX, minY, maxX, maxY int, sY, sCb, sCr []uint8, yStride, cStride int)  uint64")
	Doc("AsmYCbCrToGray is a forward DCT transform for []float32")

	// Load 128 x4
	static := make([]reg.VecVirtual, 10)
	for i := 0; i < len(static); i++ {
		static[i] = XMM()
		VPBROADCASTD(pixelsToGray4Values.Offset(i*4), static[i])
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

	y := GP64()
	x := GP64()
	idx := GP64()
	XORQ(y, y)     // set y to zero
	XORQ(x, x)     // set x to zero
	XORQ(idx, idx) // set idx to zero

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
	TEXT("AsmYCbCrToGray8", NOSPLIT, "func(pixels []float32, minX, minY, maxX, maxY int, sY, sCb, sCr []uint8, yStride, cStride int)  uint64")
	Doc("AsmYCbCrToGray8 is a forward DCT transform for []float32")

	// Load 128 x4
	static := make([]reg.VecVirtual, 10)
	for i := 0; i < len(static); i++ {
		static[i] = YMM()
		VPBROADCASTD(pixelsToGray8Values.Offset(i*4), static[i])
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
