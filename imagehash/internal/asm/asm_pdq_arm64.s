//go:build arm64

#include "textflag.h"

// PDQ luminance kernels. Each lane computes the exact scalar operation
// sequence from ycbcrToLuma (integer JFIF ratios, arithmetic shift, saturate
// to [0,255], float BT.601 weights with separate multiply and add), so output
// is bit-identical to the portable Go code on every input.
//
// Both kernels process 8 pixels per iteration. Length and alignment
// preconditions are enforced by the Go wrapper; the kernels assume n > 0 and
// n a multiple of 8.

// asmPDQLuma444Row computes luminance for n 1:1-mapped YCbCr pixels.
//
// ABI0 layout:
//   out []float32  out+0(FP)
//   y   []uint8    y+24(FP)
//   cb  []uint8    cb+48(FP)
//   cr  []uint8    cr+72(FP)
TEXT ·asmPDQLuma444Row(SB), NOSPLIT, $0-96
	MOVD	out_base+0(FP), R0
	MOVD	y_base+24(FP), R1
	MOVD	cb_base+48(FP), R2
	MOVD	cr_base+72(FP), R3
	MOVD	out_len+8(FP), R4

	MOVD	$pdqLumaConst<>(SB), R10
	VLD1R	(R10), [V16.S4]       // 128
	ADD	$4, R10, R10
	VLD1R	(R10), [V17.S4]       // 0x10101
	ADD	$4, R10, R10
	VLD1R	(R10), [V18.S4]       // 91881
	ADD	$4, R10, R10
	VLD1R	(R10), [V19.S4]       // 46802
	ADD	$4, R10, R10
	VLD1R	(R10), [V20.S4]       // 22554
	ADD	$4, R10, R10
	VLD1R	(R10), [V21.S4]       // 116130
	ADD	$4, R10, R10
	VLD1R	(R10), [V22.S4]       // 0
	ADD	$4, R10, R10
	VLD1R	(R10), [V23.S4]       // 255
	ADD	$4, R10, R10
	VLD1R	(R10), [V24.S4]       // 0.299
	ADD	$4, R10, R10
	VLD1R	(R10), [V25.S4]       // 0.587
	ADD	$4, R10, R10
	VLD1R	(R10), [V26.S4]       // 0.114

loop444:
	CMP	$0, R4
	BEQ	done444
	VLD1	(R1), [V0.B8]
	VLD1	(R2), [V1.B8]
	VLD1	(R3), [V2.B8]
	VUXTL	V0.B8, V3.H8
	VUXTL	V1.B8, V4.H8
	VUXTL	V2.B8, V5.H8
	VUXTL	V3.H4, V6.S4          // Y low
	VUXTL	V4.H4, V7.S4          // Cb low
	VUXTL	V5.H4, V8.S4          // Cr low
	VUXTL2	V3.H8, V9.S4          // Y high
	VUXTL2	V4.H8, V10.S4         // Cb high
	VUXTL2	V5.H8, V11.S4         // Cr high

	// Low four pixels.
	VMUL	V17.S4, V6.S4, V0.S4    // yy
	VSUB	V16.S4, V7.S4, V1.S4    // cb1
	VSUB	V16.S4, V8.S4, V2.S4    // cr1
	VMUL	V18.S4, V2.S4, V3.S4
	VADD	V0.S4, V3.S4, V3.S4     // r = yy + 91881*cr
	VMUL	V20.S4, V1.S4, V4.S4
	VSUB	V4.S4, V0.S4, V4.S4
	VMUL	V19.S4, V2.S4, V5.S4
	VSUB	V5.S4, V4.S4, V4.S4     // g = yy - 22554*cb - 46802*cr
	VMUL	V21.S4, V1.S4, V5.S4
	VADD	V0.S4, V5.S4, V5.S4     // b = yy + 116130*cb
	VSSHR	$16, V3.S4, V3.S4
	VSSHR	$16, V4.S4, V4.S4
	VSSHR	$16, V5.S4, V5.S4
	VSMAX	V22.S4, V3.S4, V3.S4
	VSMAX	V22.S4, V4.S4, V4.S4
	VSMAX	V22.S4, V5.S4, V5.S4
	VSMIN	V23.S4, V3.S4, V3.S4
	VSMIN	V23.S4, V4.S4, V4.S4
	VSMIN	V23.S4, V5.S4, V5.S4
	VSCVTF	V3.S4, V3.S4
	VSCVTF	V4.S4, V4.S4
	VSCVTF	V5.S4, V5.S4
	VFMUL	V24.S4, V3.S4, V3.S4
	VFMLA	V25.S4, V4.S4, V3.S4
	VFMLA	V26.S4, V5.S4, V3.S4
	VMOV	V3.B16, V12.B16

	// High four pixels.
	VMUL	V17.S4, V9.S4, V0.S4
	VSUB	V16.S4, V10.S4, V1.S4
	VSUB	V16.S4, V11.S4, V2.S4
	VMUL	V18.S4, V2.S4, V3.S4
	VADD	V0.S4, V3.S4, V3.S4
	VMUL	V20.S4, V1.S4, V4.S4
	VSUB	V4.S4, V0.S4, V4.S4
	VMUL	V19.S4, V2.S4, V5.S4
	VSUB	V5.S4, V4.S4, V4.S4
	VMUL	V21.S4, V1.S4, V5.S4
	VADD	V0.S4, V5.S4, V5.S4
	VSSHR	$16, V3.S4, V3.S4
	VSSHR	$16, V4.S4, V4.S4
	VSSHR	$16, V5.S4, V5.S4
	VSMAX	V22.S4, V3.S4, V3.S4
	VSMAX	V22.S4, V4.S4, V4.S4
	VSMAX	V22.S4, V5.S4, V5.S4
	VSMIN	V23.S4, V3.S4, V3.S4
	VSMIN	V23.S4, V4.S4, V4.S4
	VSMIN	V23.S4, V5.S4, V5.S4
	VSCVTF	V3.S4, V3.S4
	VSCVTF	V4.S4, V4.S4
	VSCVTF	V5.S4, V5.S4
	VFMUL	V24.S4, V3.S4, V3.S4
	VFMLA	V25.S4, V4.S4, V3.S4
	VFMLA	V26.S4, V5.S4, V3.S4
	VMOV	V3.B16, V13.B16

	VST1	[V12.S4], (R0)
	ADD	$16, R0, R5
	VST1	[V13.S4], (R5)
	SUB	$16, R5, R0

	ADD	$8, R1, R1
	ADD	$8, R2, R2
	ADD	$8, R3, R3
	ADD	$32, R0, R0
	SUB	$8, R4, R4
	B	loop444
done444:
	RET

// asmPDQLuma420Row computes luminance for n Y pixels sharing n/2 Cb and n/2
// Cr pixels (horizontal 2:1 subsampling). Chroma bytes are duplicated to
// pairs with VTBL, matching (minX+x)/2 indexing for even origins.
//
// ABI0 layout is identical to asmPDQLuma444Row, except cb and cr hold n/2
// bytes each.
TEXT ·asmPDQLuma420Row(SB), NOSPLIT, $0-96
	MOVD	out_base+0(FP), R0
	MOVD	y_base+24(FP), R1
	MOVD	cb_base+48(FP), R2
	MOVD	cr_base+72(FP), R3
	MOVD	out_len+8(FP), R4

	MOVD	$pdqLumaConst<>(SB), R10
	VLD1R	(R10), [V16.S4]
	ADD	$4, R10, R10
	VLD1R	(R10), [V17.S4]
	ADD	$4, R10, R10
	VLD1R	(R10), [V18.S4]
	ADD	$4, R10, R10
	VLD1R	(R10), [V19.S4]
	ADD	$4, R10, R10
	VLD1R	(R10), [V20.S4]
	ADD	$4, R10, R10
	VLD1R	(R10), [V21.S4]
	ADD	$4, R10, R10
	VLD1R	(R10), [V22.S4]
	ADD	$4, R10, R10
	VLD1R	(R10), [V23.S4]
	ADD	$4, R10, R10
	VLD1R	(R10), [V24.S4]
	ADD	$4, R10, R10
	VLD1R	(R10), [V25.S4]
	ADD	$4, R10, R10
	VLD1R	(R10), [V26.S4]
	MOVD	$pdqDupIdx<>(SB), R11
	VLD1	(R11), [V27.B8]

loop420:
	CMP	$0, R4
	BEQ	done420
	VLD1	(R1), [V0.B8]
	// Single-lane loads read exactly the 4 chroma bytes consumed below;
	// VTBL indexes only bytes 0-3, so stale upper lanes are never read.
	VLD1	(R2), V1.S[0]
	VLD1	(R3), V2.S[0]
	// Duplicate chroma bytes to pairs. Note Go's VTBL operand order is
	// index, table, dest.
	VTBL	V27.B8, [V1.B16], V4.B8   // cb pairs
	VTBL	V27.B8, [V2.B16], V5.B8   // cr pairs
	VUXTL	V0.B8, V3.H8
	VUXTL	V4.B8, V6.H8
	VUXTL	V5.B8, V14.H8
	VUXTL	V3.H4, V7.S4          // Y low
	VUXTL	V6.H4, V8.S4          // Cb low
	VUXTL	V14.H4, V9.S4         // Cr low
	VUXTL2	V3.H8, V10.S4         // Y high
	VUXTL2	V6.H8, V11.S4         // Cb high
	VUXTL2	V14.H8, V13.S4        // Cr high (V13 reused for output below)

	VMUL	V17.S4, V7.S4, V0.S4
	VSUB	V16.S4, V8.S4, V1.S4
	VSUB	V16.S4, V9.S4, V2.S4
	VMUL	V18.S4, V2.S4, V3.S4
	VADD	V0.S4, V3.S4, V3.S4
	VMUL	V20.S4, V1.S4, V4.S4
	VSUB	V4.S4, V0.S4, V4.S4
	VMUL	V19.S4, V2.S4, V5.S4
	VSUB	V5.S4, V4.S4, V4.S4
	VMUL	V21.S4, V1.S4, V5.S4
	VADD	V0.S4, V5.S4, V5.S4
	VSSHR	$16, V3.S4, V3.S4
	VSSHR	$16, V4.S4, V4.S4
	VSSHR	$16, V5.S4, V5.S4
	VSMAX	V22.S4, V3.S4, V3.S4
	VSMAX	V22.S4, V4.S4, V4.S4
	VSMAX	V22.S4, V5.S4, V5.S4
	VSMIN	V23.S4, V3.S4, V3.S4
	VSMIN	V23.S4, V4.S4, V4.S4
	VSMIN	V23.S4, V5.S4, V5.S4
	VSCVTF	V3.S4, V3.S4
	VSCVTF	V4.S4, V4.S4
	VSCVTF	V5.S4, V5.S4
	VFMUL	V24.S4, V3.S4, V3.S4
	VFMLA	V25.S4, V4.S4, V3.S4
	VFMLA	V26.S4, V5.S4, V3.S4
	VMOV	V3.B16, V12.B16

	VMUL	V17.S4, V10.S4, V0.S4
	VSUB	V16.S4, V11.S4, V1.S4
	VSUB	V16.S4, V13.S4, V2.S4
	VMUL	V18.S4, V2.S4, V3.S4
	VADD	V0.S4, V3.S4, V3.S4
	VMUL	V20.S4, V1.S4, V4.S4
	VSUB	V4.S4, V0.S4, V4.S4
	VMUL	V19.S4, V2.S4, V5.S4
	VSUB	V5.S4, V4.S4, V4.S4
	VMUL	V21.S4, V1.S4, V5.S4
	VADD	V0.S4, V5.S4, V5.S4
	VSSHR	$16, V3.S4, V3.S4
	VSSHR	$16, V4.S4, V4.S4
	VSSHR	$16, V5.S4, V5.S4
	VSMAX	V22.S4, V3.S4, V3.S4
	VSMAX	V22.S4, V4.S4, V4.S4
	VSMAX	V22.S4, V5.S4, V5.S4
	VSMIN	V23.S4, V3.S4, V3.S4
	VSMIN	V23.S4, V4.S4, V4.S4
	VSMIN	V23.S4, V5.S4, V5.S4
	VSCVTF	V3.S4, V3.S4
	VSCVTF	V4.S4, V4.S4
	VSCVTF	V5.S4, V5.S4
	VFMUL	V24.S4, V3.S4, V3.S4
	VFMLA	V25.S4, V4.S4, V3.S4
	VFMLA	V26.S4, V5.S4, V3.S4
	VMOV	V3.B16, V13.B16

	VST1	[V12.S4], (R0)
	ADD	$16, R0, R5
	VST1	[V13.S4], (R5)
	SUB	$16, R5, R0

	ADD	$8, R1, R1
	ADD	$4, R2, R2
	ADD	$4, R3, R3
	ADD	$32, R0, R0
	SUB	$8, R4, R4
	B	loop420
done420:
	RET

DATA pdqLumaConst<>+0(SB)/4, $128
DATA pdqLumaConst<>+4(SB)/4, $0x10101
DATA pdqLumaConst<>+8(SB)/4, $91881
DATA pdqLumaConst<>+12(SB)/4, $46802
DATA pdqLumaConst<>+16(SB)/4, $22554
DATA pdqLumaConst<>+20(SB)/4, $116130
DATA pdqLumaConst<>+24(SB)/4, $0
DATA pdqLumaConst<>+28(SB)/4, $255
DATA pdqLumaConst<>+32(SB)/4, $0x3E991687 // F32(0.299)
DATA pdqLumaConst<>+36(SB)/4, $0x3F1645A2 // F32(0.587)
DATA pdqLumaConst<>+40(SB)/4, $0x3DE978D5 // F32(0.114)
GLOBL pdqLumaConst<>(SB), RODATA|NOPTR, $44

// Pair-duplication indices: byte i of each pair maps to chroma byte i/2.
DATA pdqDupIdx<>+0(SB)/1, $0
DATA pdqDupIdx<>+1(SB)/1, $0
DATA pdqDupIdx<>+2(SB)/1, $1
DATA pdqDupIdx<>+3(SB)/1, $1
DATA pdqDupIdx<>+4(SB)/1, $2
DATA pdqDupIdx<>+5(SB)/1, $2
DATA pdqDupIdx<>+6(SB)/1, $3
DATA pdqDupIdx<>+7(SB)/1, $3
GLOBL pdqDupIdx<>(SB), RODATA|NOPTR, $8
