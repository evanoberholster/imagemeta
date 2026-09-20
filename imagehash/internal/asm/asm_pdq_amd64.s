//go:build amd64

// Hand-written AVX2 PDQ luminance kernels (the avo generator in asm.go covers
// the pHash kernels; these follow the same ABI and style).
//
// Each lane computes the exact scalar operation sequence from ycbcrToLuma
// (integer JFIF ratios, arithmetic shift, saturate to [0,255], float BT.601
// weights with SEPARATE multiply and add), so output is bit-identical to the
// portable Go code on amd64. Note the scalar uses separate MULSS/ADDSS on
// amd64 (unlike arm64, where the compiler contracts FMADDS), hence no FMA is
// used here either.
//
// Both kernels process 8 pixels per iteration. Length and alignment
// preconditions are enforced by the Go wrapper; the kernels assume n > 0 and
// n a multiple of 8.
//
// Register layout: Y0-Y5 temporaries,
// Y6=128, Y7=0x10101, Y8=91881, Y9=46802, Y10=22554, Y11=116130, Y12=255,
// Y13=0.299, Y14=0.587, Y15=0.114.
//
// ABI0 layout:
//   out []float32  out+0(FP)
//   y   []uint8    y+24(FP)
//   cb  []uint8    cb+48(FP)
//   cr  []uint8    cr+72(FP)

#include "textflag.h"

// asmPDQLuma444Row computes luminance for n 1:1-mapped YCbCr pixels.
TEXT ·asmPDQLuma444Row(SB), NOSPLIT|NOPTR, $0-96
	MOVQ	out_base+0(FP), DI
	MOVQ	y_base+24(FP), SI
	MOVQ	cb_base+48(FP), DX
	MOVQ	cr_base+72(FP), CX
	MOVQ	out_len+8(FP), R8

	VPBROADCASTD	pdqLumaInt<>+0(SB), Y6
	VPBROADCASTD	pdqLumaInt<>+4(SB), Y7
	VPBROADCASTD	pdqLumaInt<>+8(SB), Y8
	VPBROADCASTD	pdqLumaInt<>+12(SB), Y9
	VPBROADCASTD	pdqLumaInt<>+16(SB), Y10
	VPBROADCASTD	pdqLumaInt<>+20(SB), Y11
	VPBROADCASTD	pdqLumaInt<>+24(SB), Y12
	VBROADCASTSS	pdqLumaFloat<>+0(SB), Y13
	VBROADCASTSS	pdqLumaFloat<>+4(SB), Y14
	VBROADCASTSS	pdqLumaFloat<>+8(SB), Y15

loop444:
	TESTQ	R8, R8
	JE	done444
	VPMOVZXBD	(SI), Y0
	VPMOVZXBD	(DX), Y1
	VPMOVZXBD	(CX), Y2
	VPMULLD	Y7, Y0, Y0            // yy
	VPSUBD	Y6, Y1, Y1            // cb1
	VPSUBD	Y6, Y2, Y2            // cr1
	VPMULLD	Y8, Y2, Y3
	VPADDD	Y0, Y3, Y3            // r
	VPMULLD	Y10, Y1, Y4
	VPSUBD	Y4, Y0, Y4
	VPMULLD	Y9, Y2, Y5
	VPSUBD	Y5, Y4, Y4            // g
	VPMULLD	Y11, Y1, Y5
	VPADDD	Y0, Y5, Y5            // b
	VPSRAD	$16, Y3, Y3
	VPSRAD	$16, Y4, Y4
	VPSRAD	$16, Y5, Y5
	VPXOR	Y0, Y0, Y0
	VPMAXSD	Y0, Y3, Y3
	VPMAXSD	Y0, Y4, Y4
	VPMAXSD	Y0, Y5, Y5
	VPMINSD	Y12, Y3, Y3
	VPMINSD	Y12, Y4, Y4
	VPMINSD	Y12, Y5, Y5
	VCVTDQ2PS	Y3, Y3
	VCVTDQ2PS	Y4, Y4
	VCVTDQ2PS	Y5, Y5
	VMULPS	Y13, Y3, Y3
	VMULPS	Y14, Y4, Y4
	VMULPS	Y15, Y5, Y5
	VADDPS	Y4, Y3, Y3
	VADDPS	Y5, Y3, Y3
	VMOVDQU	Y3, (DI)

	ADDQ	$8, SI
	ADDQ	$8, DX
	ADDQ	$8, CX
	ADDQ	$32, DI
	SUBQ	$8, R8
	JMP	loop444
done444:
	VZEROUPPER
	RET

// asmPDQLuma420Row computes luminance for n Y pixels sharing n/2 Cb and n/2
// Cr pixels (horizontal 2:1 subsampling). Chroma dwords are broadcast and
// shuffled to pairs, matching (minX+x)/2 indexing for even origins.
TEXT ·asmPDQLuma420Row(SB), NOSPLIT|NOPTR, $0-96
	MOVQ	out_base+0(FP), DI
	MOVQ	y_base+24(FP), SI
	MOVQ	cb_base+48(FP), DX
	MOVQ	cr_base+72(FP), CX
	MOVQ	out_len+8(FP), R8

	VPBROADCASTD	pdqLumaInt<>+0(SB), Y6
	VPBROADCASTD	pdqLumaInt<>+4(SB), Y7
	VPBROADCASTD	pdqLumaInt<>+8(SB), Y8
	VPBROADCASTD	pdqLumaInt<>+12(SB), Y9
	VPBROADCASTD	pdqLumaInt<>+16(SB), Y10
	VPBROADCASTD	pdqLumaInt<>+20(SB), Y11
	VPBROADCASTD	pdqLumaInt<>+24(SB), Y12
	VBROADCASTSS	pdqLumaFloat<>+0(SB), Y13
	VBROADCASTSS	pdqLumaFloat<>+4(SB), Y14
	VBROADCASTSS	pdqLumaFloat<>+8(SB), Y15

loop420:
	TESTQ	R8, R8
	JE	done420
	VPMOVZXBD	(SI), Y0
	VPBROADCASTD	(DX), X1
	VPBROADCASTD	(CX), X2
	VPSHUFB	pdqDupMask<>(SB), X1, X1
	VPSHUFB	pdqDupMask<>(SB), X2, X2
	VPMOVZXBD	X1, Y1
	VPMOVZXBD	X2, Y2
	VPMULLD	Y7, Y0, Y0
	VPSUBD	Y6, Y1, Y1
	VPSUBD	Y6, Y2, Y2
	VPMULLD	Y8, Y2, Y3
	VPADDD	Y0, Y3, Y3
	VPMULLD	Y10, Y1, Y4
	VPSUBD	Y4, Y0, Y4
	VPMULLD	Y9, Y2, Y5
	VPSUBD	Y5, Y4, Y4
	VPMULLD	Y11, Y1, Y5
	VPADDD	Y0, Y5, Y5
	VPSRAD	$16, Y3, Y3
	VPSRAD	$16, Y4, Y4
	VPSRAD	$16, Y5, Y5
	VPXOR	Y0, Y0, Y0
	VPMAXSD	Y0, Y3, Y3
	VPMAXSD	Y0, Y4, Y4
	VPMAXSD	Y0, Y5, Y5
	VPMINSD	Y12, Y3, Y3
	VPMINSD	Y12, Y4, Y4
	VPMINSD	Y12, Y5, Y5
	VCVTDQ2PS	Y3, Y3
	VCVTDQ2PS	Y4, Y4
	VCVTDQ2PS	Y5, Y5
	VMULPS	Y13, Y3, Y3
	VMULPS	Y14, Y4, Y4
	VMULPS	Y15, Y5, Y5
	VADDPS	Y4, Y3, Y3
	VADDPS	Y5, Y3, Y3
	VMOVDQU	Y3, (DI)

	ADDQ	$8, SI
	ADDQ	$4, DX
	ADDQ	$4, CX
	ADDQ	$32, DI
	SUBQ	$8, R8
	JMP	loop420
done420:
	VZEROUPPER
	RET

DATA pdqLumaInt<>+0(SB)/4, $128
DATA pdqLumaInt<>+4(SB)/4, $0x10101
DATA pdqLumaInt<>+8(SB)/4, $91881
DATA pdqLumaInt<>+12(SB)/4, $46802
DATA pdqLumaInt<>+16(SB)/4, $22554
DATA pdqLumaInt<>+20(SB)/4, $116130
DATA pdqLumaInt<>+24(SB)/4, $255
GLOBL pdqLumaInt<>(SB), RODATA|NOPTR, $28

DATA pdqLumaFloat<>+0(SB)/4, $0x3E991687 // F32(0.299)
DATA pdqLumaFloat<>+4(SB)/4, $0x3F1645A2 // F32(0.587)
DATA pdqLumaFloat<>+8(SB)/4, $0x3DE978D5 // F32(0.114)
GLOBL pdqLumaFloat<>(SB), RODATA|NOPTR, $12

// Pair-duplication shuffle mask: byte i of each 8-byte half maps to chroma
// byte i/2.
DATA pdqDupMask<>+0(SB)/1, $0
DATA pdqDupMask<>+1(SB)/1, $0
DATA pdqDupMask<>+2(SB)/1, $1
DATA pdqDupMask<>+3(SB)/1, $1
DATA pdqDupMask<>+4(SB)/1, $2
DATA pdqDupMask<>+5(SB)/1, $2
DATA pdqDupMask<>+6(SB)/1, $3
DATA pdqDupMask<>+7(SB)/1, $3
DATA pdqDupMask<>+8(SB)/1, $0
DATA pdqDupMask<>+9(SB)/1, $0
DATA pdqDupMask<>+10(SB)/1, $1
DATA pdqDupMask<>+11(SB)/1, $1
DATA pdqDupMask<>+12(SB)/1, $2
DATA pdqDupMask<>+13(SB)/1, $2
DATA pdqDupMask<>+14(SB)/1, $3
DATA pdqDupMask<>+15(SB)/1, $3
GLOBL pdqDupMask<>(SB), RODATA|NOPTR, $16
