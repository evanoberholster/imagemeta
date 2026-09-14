//go:build arm64

#include "textflag.h"

// asmDCT2DHash64 is an ARM64 NEON implementation of the DCT-II 2D transform
// used by PHash64. It is a direct port of transforms32.DCT2DHash64 and uses
// the same iterative Lee algorithm, producing identical results.
//
// It transforms all 64 rows in place, then transforms the first 8 columns and
// keeps the top-left 8x8 block as [64]float32.
//
// ABI0 layout:
//   input_base  input+0(FP), input_len+8(FP), input_cap+16(FP)
//   return      ret_0+24(FP) ... [64]float32
TEXT ·asmDCT2DHash64(SB), NOSPLIT, $544-280
	MOVD	input_base+0(FP), R0
	MOVD	$ret_0+24(FP), R1
	MOVD	$scr-272(SP), R20
	MOVD	$col-528(SP), R21

	// Reverse-lane permutation table for 4 float32 lanes, shared by all DCTs.
	MOVD	$revtbl<>(SB), R25
	VLD1	(R25), [V7.B16]

	// Iterate 72 transforms: 0..63 are rows, 64..71 are columns.
	MOVD	$0, R19
mainloop:
	CMP	$64, R19
	BHS	do_col

	// Row setup: R3 = input + index*256
	LSL	$8, R19, R10
	ADD	R0, R10, R3
	MOVD	R20, R4
	B	run_dct

do_col:
	// Column index = index - 64
	SUB	$64, R19, R23
	LSL	$2, R23, R23
	MOVD	$0, R7
gatherloop:
	LSL	$8, R7, R10
	ADD	R0, R10, R10
	ADD	R23, R10, R10
	FMOVS	(R10), F0
	LSL	$2, R7, R13
	ADD	R21, R13, R13
	FMOVS	F0, (R13)
	ADD	$1, R7, R7
	CMP	$64, R7
	BNE	gatherloop
	MOVD	R21, R3
	MOVD	R20, R4

run_dct:
	// ---------------------------------------------------------------
	// In-place 64-point DCT-II on R3 using R4 as scratch.
	// Iterative Lee algorithm: 6 split levels down, 6 merge levels up.
	// ---------------------------------------------------------------
	MOVD	$256, R5
	MOVD	$128, R8
	MOVD	$dcttbl<>(SB), R6
	MOVD	$128, R16
	MOVD	$6, R17
dloop:
	CMP	$16, R5
	BEQ	dsplit_m4
	CMP	$8, R5
	BEQ	dsplit_m2
	MOVD	$0, R9
dsplit_block:
	CMP	$256, R9
	BHS	dsplit_blockdone
	ADD	R3, R9, R23
	ADD	R4, R9, R24
	MOVD	$0, R7
dsplit_v:
	ADD	$16, R7, R25
	CMP	R8, R25
	BHI	dsplit_s
	ADD	R23, R7, R10
	VLD1	(R10), [V0.S4]
	SUB	R7, R5, R11
	SUB	$16, R11, R11
	ADD	R23, R11, R11
	VLD1	(R11), [V1.S4]
	VTBL	V7.B16, [V1.B16], V1.B16
	VFADD	V1.S4, V0.S4, V2.S4
	VFSUB	V1.S4, V0.S4, V3.S4
	ADD	R6, R7, R12
	VLD1	(R12), [V4.S4]
	VFDIV	V4.S4, V3.S4, V3.S4
	ADD	R24, R7, R12
	VST1	[V2.S4], (R12)
	ADD	R8, R7, R13
	ADD	R24, R13, R13
	VST1	[V3.S4], (R13)
	ADD	$16, R7, R7
	B	dsplit_v
dsplit_s:
	CMP	R8, R7
	BHS	dsplit_blocknext
	ADD	R23, R7, R10
	FMOVS	(R10), F0
	SUB	R7, R5, R11
	SUB	$4, R11, R11
	ADD	R23, R11, R11
	FMOVS	(R11), F1
	FADDS	F1, F0, F2
	FSUBS	F1, F0, F3
	ADD	R6, R7, R12
	FMOVS	(R12), F4
	FDIVS	F4, F3, F3
	ADD	R24, R7, R12
	FMOVS	F2, (R12)
	ADD	R8, R7, R13
	ADD	R24, R13, R13
	FMOVS	F3, (R13)
	ADD	$4, R7, R7
	B	dsplit_s
dsplit_blocknext:
	ADD	R5, R9, R9
	B	dsplit_block

// M=4 split, four blocks at a time. dcttbl+240 holds the two M=4 twiddles.
dsplit_m4:
	MOVD	$dcttbl<>(SB), R15
	ADD	$240, R15, R15
	VLD1R	(R15), [V20.S4]
	ADD	$4, R15, R15
	VLD1R	(R15), [V21.S4]
	MOVD	$0, R9
dsplit_m4_loop:
	ADD	R3, R9, R23
	ADD	R4, R9, R24
	VLD4	(R23), [V0.S4, V1.S4, V2.S4, V3.S4]
	VFADD	V3.S4, V0.S4, V16.S4
	VFADD	V2.S4, V1.S4, V17.S4
	VFSUB	V3.S4, V0.S4, V18.S4
	VFSUB	V2.S4, V1.S4, V19.S4
	VFDIV	V20.S4, V18.S4, V18.S4
	VFDIV	V21.S4, V19.S4, V19.S4
	VST4	[V16.S4, V17.S4, V18.S4, V19.S4], (R24)
	ADD	$64, R9, R9
	CMP	$256, R9
	BNE	dsplit_m4_loop
	B	dsplit_blockdone

// M=2 split, four blocks at a time. dcttbl+248 holds the M=2 twiddle.
dsplit_m2:
	MOVD	$dcttbl<>(SB), R15
	ADD	$248, R15, R15
	VLD1R	(R15), [V22.S4]
	MOVD	$0, R9
dsplit_m2_loop:
	ADD	R3, R9, R23
	ADD	R4, R9, R24
	VLD2	(R23), [V0.S4, V1.S4]
	VFADD	V1.S4, V0.S4, V16.S4
	VFSUB	V1.S4, V0.S4, V17.S4
	VFDIV	V22.S4, V17.S4, V17.S4
	VST2	[V16.S4, V17.S4], (R24)
	ADD	$32, R9, R9
	CMP	$256, R9
	BNE	dsplit_m2_loop
	B	dsplit_blockdone

dsplit_blockdone:
	ADD	R16, R6, R6
	LSR	$1, R5, R5
	LSR	$1, R8, R8
	LSR	$1, R16, R16
	MOVD	R3, R25
	MOVD	R4, R3
	MOVD	R25, R4
	SUB	$1, R17, R17
	CBNZ	R17, dloop

	// Upward recombination pass.
	MOVD	$8, R5
	MOVD	$4, R8
	MOVD	$6, R17
uloop:
	CMP	$8, R5
	BEQ	umerge_m2
	CMP	$16, R5
	BEQ	umerge_m4
	MOVD	$0, R9
umerge_block:
	CMP	$256, R9
	BHS	umerge_done
	ADD	R3, R9, R23
	ADD	R4, R9, R24
	SUB	$4, R8, R14
	MOVD	$0, R7
umerge_v:
	ADD	$16, R7, R25
	CMP	R14, R25
	BHI	umerge_s
	ADD	R23, R7, R10
	VLD1	(R10), [V0.S4]
	ADD	R8, R7, R11
	ADD	R23, R11, R11
	VLD1	(R11), [V1.S4]
	ADD	$4, R11, R15
	VLD1	(R15), [V2.S4]
	VFADD	V2.S4, V1.S4, V1.S4
	VZIP1	V1.S4, V0.S4, V3.S4
	VZIP2	V1.S4, V0.S4, V4.S4
	LSL	$1, R7, R12
	ADD	R24, R12, R12
	VST1	[V3.S4], (R12)
	ADD	$16, R12, R15
	VST1	[V4.S4], (R15)
	ADD	$16, R7, R7
	B	umerge_v
umerge_s:
	CMP	R14, R7
	BHS	umerge_tail
	ADD	R23, R7, R10
	FMOVS	(R10), F0
	ADD	R8, R7, R11
	ADD	R23, R11, R11
	FMOVS	(R11), F1
	FMOVS	4(R11), F2
	FADDS	F2, F1, F3
	LSL	$1, R7, R12
	ADD	R24, R12, R12
	FMOVS	F0, (R12)
	FMOVS	F3, 4(R12)
	ADD	$4, R7, R7
	B	umerge_s
umerge_tail:
	SUB	$8, R5, R10
	ADD	R24, R10, R10
	SUB	$4, R8, R11
	ADD	R23, R11, R11
	FMOVS	(R11), F0
	FMOVS	F0, (R10)
	ADD	R23, R5, R12
	SUB	$4, R12, R12
	FMOVS	(R12), F0
	ADD	R24, R5, R13
	SUB	$4, R13, R13
	FMOVS	F0, (R13)
	ADD	R5, R9, R9
	B	umerge_block

// M=2 merge is an identity copy, done 4 floats at a time.
umerge_m2:
	MOVD	$0, R9
umerge_m2_loop:
	ADD	R3, R9, R10
	ADD	R4, R9, R11
	VLD1	(R10), [V0.S4]
	VST1	[V0.S4], (R11)
	ADD	$16, R9, R9
	CMP	$256, R9
	BNE	umerge_m2_loop
	B	umerge_done

// M=4 merge, four blocks at a time.
umerge_m4:
	MOVD	$0, R9
umerge_m4_loop:
	ADD	R3, R9, R23
	ADD	R4, R9, R24
	VLD4	(R23), [V0.S4, V1.S4, V2.S4, V3.S4]
	VFADD	V3.S4, V2.S4, V17.S4
	VMOV	V0.B16, V16.B16
	VMOV	V1.B16, V18.B16
	VMOV	V3.B16, V19.B16
	VST4	[V16.S4, V17.S4, V18.S4, V19.S4], (R24)
	ADD	$64, R9, R9
	CMP	$256, R9
	BNE	umerge_m4_loop
	B	umerge_done

umerge_done:
	MOVD	R3, R25
	MOVD	R4, R3
	MOVD	R25, R4
	LSL	$1, R5, R5
	LSL	$1, R8, R8
	SUB	$1, R17, R17
	CBNZ	R17, uloop

	// ---------------------------------------------------------------
	CMP	$64, R19
	BHS	col_teardown
	B	next

col_teardown:
	SUB	$64, R19, R23
	LSL	$2, R23, R23
	MOVD	$0, R7
scatterloop:
	LSL	$2, R7, R13
	ADD	R21, R13, R13
	FMOVS	(R13), F0
	LSL	$5, R7, R14
	ADD	R23, R14, R14
	ADD	R1, R14, R14
	FMOVS	F0, (R14)
	ADD	$1, R7, R7
	CMP	$8, R7
	BNE	scatterloop

next:
	ADD	$1, R19, R19
	CMP	$72, R19
	BNE	mainloop
	RET

DATA revtbl<>+0(SB)/1, $12
DATA revtbl<>+1(SB)/1, $13
DATA revtbl<>+2(SB)/1, $14
DATA revtbl<>+3(SB)/1, $15
DATA revtbl<>+4(SB)/1, $8
DATA revtbl<>+5(SB)/1, $9
DATA revtbl<>+6(SB)/1, $10
DATA revtbl<>+7(SB)/1, $11
DATA revtbl<>+8(SB)/1, $4
DATA revtbl<>+9(SB)/1, $5
DATA revtbl<>+10(SB)/1, $6
DATA revtbl<>+11(SB)/1, $7
DATA revtbl<>+12(SB)/1, $0
DATA revtbl<>+13(SB)/1, $1
DATA revtbl<>+14(SB)/1, $2
DATA revtbl<>+15(SB)/1, $3
GLOBL revtbl<>(SB), RODATA|NOPTR, $16

// DCT-II twiddle divisors 2*cos((i+0.5)*pi/M) for M = 64,32,16,8,4,2.
DATA dcttbl<>+0(SB)/4, $1.99939764
DATA dcttbl<>+4(SB)/4, $1.99458086
DATA dcttbl<>+8(SB)/4, $1.98495913
DATA dcttbl<>+12(SB)/4, $1.97055531
DATA dcttbl<>+16(SB)/4, $1.95140421
DATA dcttbl<>+20(SB)/4, $1.9275521
DATA dcttbl<>+24(SB)/4, $1.89905632
DATA dcttbl<>+28(SB)/4, $1.86598563
DATA dcttbl<>+32(SB)/4, $1.82841957
DATA dcttbl<>+36(SB)/4, $1.7864486
DATA dcttbl<>+40(SB)/4, $1.74017394
DATA dcttbl<>+44(SB)/4, $1.68970716
DATA dcttbl<>+48(SB)/4, $1.63516963
DATA dcttbl<>+52(SB)/4, $1.57669282
DATA dcttbl<>+56(SB)/4, $1.51441765
DATA dcttbl<>+60(SB)/4, $1.4484942
DATA dcttbl<>+64(SB)/4, $1.37908113
DATA dcttbl<>+68(SB)/4, $1.3063457
DATA dcttbl<>+72(SB)/4, $1.23046315
DATA dcttbl<>+76(SB)/4, $1.15161633
DATA dcttbl<>+80(SB)/4, $1.06999528
DATA dcttbl<>+84(SB)/4, $0.985796392
DATA dcttbl<>+88(SB)/4, $0.899222672
DATA dcttbl<>+92(SB)/4, $0.810482621
DATA dcttbl<>+96(SB)/4, $0.719790101
DATA dcttbl<>+100(SB)/4, $0.627363503
DATA dcttbl<>+104(SB)/4, $0.53342551
DATA dcttbl<>+108(SB)/4, $0.438202471
DATA dcttbl<>+112(SB)/4, $0.341923773
DATA dcttbl<>+116(SB)/4, $0.244821355
DATA dcttbl<>+120(SB)/4, $0.147129133
DATA dcttbl<>+124(SB)/4, $0.049082458
DATA dcttbl<>+128(SB)/4, $1.9975909
DATA dcttbl<>+132(SB)/4, $1.97835302
DATA dcttbl<>+136(SB)/4, $1.94006252
DATA dcttbl<>+140(SB)/4, $1.88308811
DATA dcttbl<>+144(SB)/4, $1.80797863
DATA dcttbl<>+148(SB)/4, $1.7154572
DATA dcttbl<>+152(SB)/4, $1.60641503
DATA dcttbl<>+156(SB)/4, $1.48190224
DATA dcttbl<>+160(SB)/4, $1.34311795
DATA dcttbl<>+164(SB)/4, $1.19139862
DATA dcttbl<>+168(SB)/4, $1.02820551
DATA dcttbl<>+172(SB)/4, $0.855110168
DATA dcttbl<>+176(SB)/4, $0.673779726
DATA dcttbl<>+180(SB)/4, $0.485960364
DATA dcttbl<>+184(SB)/4, $0.293460935
DATA dcttbl<>+188(SB)/4, $0.0981353521
DATA dcttbl<>+192(SB)/4, $1.99036944
DATA dcttbl<>+196(SB)/4, $1.91388071
DATA dcttbl<>+200(SB)/4, $1.76384258
DATA dcttbl<>+204(SB)/4, $1.54602087
DATA dcttbl<>+208(SB)/4, $1.26878655
DATA dcttbl<>+212(SB)/4, $0.942793489
DATA dcttbl<>+216(SB)/4, $0.580569327
DATA dcttbl<>+220(SB)/4, $0.196034282
DATA dcttbl<>+224(SB)/4, $1.9615705
DATA dcttbl<>+228(SB)/4, $1.66293919
DATA dcttbl<>+232(SB)/4, $1.11114049
DATA dcttbl<>+236(SB)/4, $0.390180647
DATA dcttbl<>+240(SB)/4, $1.84775901
DATA dcttbl<>+244(SB)/4, $0.765366852
DATA dcttbl<>+248(SB)/4, $1.41421354
GLOBL dcttbl<>(SB), RODATA|NOPTR, $252
