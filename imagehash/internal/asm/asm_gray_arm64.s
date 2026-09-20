//go:build arm64

#include "textflag.h"

// asmYCbCrToGray converts an *image.YCbCr with 4:4:4 chroma to grayscale
// float32 pixels. The 16-bit YCbCr to RGB conversion is done in 32-bit integer
// SIMD, shifted right by 8, then reduced to grayscale with the integer
// luminosity weights (299*r + 587*g + 114*b) / 1000. Integer-only arithmetic
// keeps the result bit-identical across architectures.
//
// The caller (yCbCrToGrayASM) guarantees SubsampleRatio == 444 and a width
// that is a multiple of 8.
//
// ABI0 layout:
//   pixels []float32  pixels+0(FP)
//   minX, minY        minX+24(FP), minY+32(FP)
//   maxX, maxY        maxX+40(FP), maxY+48(FP)
//   sY, sCb, sCr      sY+56(FP), sCb+80(FP), sCr+104(FP)
//   yStride, cStride  yStride+128(FP), cStride+136(FP)
TEXT ·asmYCbCrToGray(SB), NOSPLIT, $0-144
	MOVD	pixels_base+0(FP), R0
	MOVD	sY_base+56(FP), R1
	MOVD	sCb_base+80(FP), R2
	MOVD	sCr_base+104(FP), R3
	MOVD	yStride+128(FP), R4
	MOVD	cStride+136(FP), R5
	MOVD	minX+24(FP), R6
	MOVD	minY+32(FP), R7
	MOVD	maxX+40(FP), R8
	MOVD	maxY+48(FP), R9

	SUB	R6, R8, R19              // width = maxX - minX
	SUB	R7, R9, R21              // height = maxY - minY
	MUL	R7, R4, R10              // minY * yStride
	ADD	R6, R10, R10
	ADD	R10, R1, R10             // yp = sY + yBase
	MUL	R7, R5, R11              // minY * cStride
	ADD	R6, R11, R11
	ADD	R11, R2, R2              // cbp = sCb + cBase
	ADD	R11, R3, R3              // crp = sCr + cBase
	LSL	$2, R19, R20             // outStride = width*4
	MOVD	R0, R12                  // op = pixels
	MOVD	$0, R8                   // y = 0

	MOVD	$ycbcrConst<>(SB), R23
	VLD1R	(R23), [V16.S4]
	ADD	$4, R23, R7
	VLD1R	(R7), [V17.S4]
	ADD	$4, R7, R7
	VLD1R	(R7), [V18.S4]
	ADD	$4, R7, R7
	VLD1R	(R7), [V19.S4]
	ADD	$4, R7, R7
	VLD1R	(R7), [V20.S4]
	ADD	$4, R7, R7
	VLD1R	(R7), [V21.S4]
	ADD	$4, R7, R7
	VLD1R	(R7), [V22.S4]
	ADD	$4, R7, R7
	VLD1R	(R7), [V23.S4]
	ADD	$4, R7, R7
	VLD1R	(R7), [V24.S4]
	ADD	$4, R7, R7
	VLD1R	(R7), [V25.S4]

yloop:
	CMP	R21, R8
	BHS	done
	MOVD	R10, R13                 // yp row
	MOVD	R2, R14                  // cbp row
	MOVD	R3, R15                  // crp row
	MOVD	R12, R16                 // op row
	MOVD	$0, R17                  // x = 0
xloop:
	ADD	$8, R17, R7
	CMP	R19, R7
	BHI	xdone

	VLD1	(R13), [V0.B8]
	VLD1	(R14), [V1.B8]
	VLD1	(R15), [V2.B8]
	VUXTL	V0.B8, V3.H8
	VUXTL	V1.B8, V4.H8
	VUXTL	V2.B8, V5.H8
	VUXTL	V3.H4, V6.S4             // Y low
	VUXTL	V4.H4, V7.S4             // Cb low
	VUXTL	V5.H4, V8.S4             // Cr low
	VUXTL2	V3.H8, V9.S4             // Y high
	VUXTL2	V4.H8, V10.S4            // Cb high
	VUXTL2	V5.H8, V11.S4            // Cr high

	// Low four pixels -> V12
	VMUL	V17.S4, V6.S4, V0.S4              // yy = Y * 0x10101
	VSUB	V16.S4, V7.S4, V1.S4              // cb = Cb - 128
	VSUB	V16.S4, V8.S4, V2.S4              // cr = Cr - 128
	VMUL	V18.S4, V2.S4, V3.S4              // 91881*cr
	VADD	V0.S4, V3.S4, V3.S4               // red = yy + 91881*cr
	VMUL	V20.S4, V1.S4, V4.S4              // 22554*cb
	VSUB	V4.S4, V0.S4, V4.S4               // yy - 22554*cb
	VMUL	V19.S4, V2.S4, V5.S4              // 46802*cr
	VSUB	V5.S4, V4.S4, V4.S4               // green = yy - 22554*cb - 46802*cr
	VMUL	V21.S4, V1.S4, V5.S4              // 116130*cb
	VADD	V0.S4, V5.S4, V5.S4               // blue = yy + 116130*cb
	VSSHR	$8, V3.S4, V3.S4
	VSSHR	$8, V4.S4, V4.S4
	VSSHR	$8, V5.S4, V5.S4
	VMUL	V22.S4, V3.S4, V3.S4             // 299 * red
	VMUL	V23.S4, V4.S4, V4.S4             // 587 * green
	VMUL	V24.S4, V5.S4, V5.S4             // 114 * blue
	VADD	V4.S4, V3.S4, V3.S4
	VADD	V5.S4, V3.S4, V3.S4
	VSCVTF	V3.S4, V3.S4
	VFDIV	V25.S4, V3.S4, V3.S4             // / 1000
	VMOV	V3.B16, V12.B16

	// High four pixels -> V13
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
	VSSHR	$8, V3.S4, V3.S4
	VSSHR	$8, V4.S4, V4.S4
	VSSHR	$8, V5.S4, V5.S4
	VMUL	V22.S4, V3.S4, V3.S4             // 299 * red
	VMUL	V23.S4, V4.S4, V4.S4             // 587 * green
	VMUL	V24.S4, V5.S4, V5.S4             // 114 * blue
	VADD	V4.S4, V3.S4, V3.S4
	VADD	V5.S4, V3.S4, V3.S4
	VSCVTF	V3.S4, V3.S4
	VFDIV	V25.S4, V3.S4, V3.S4             // / 1000
	VMOV	V3.B16, V13.B16

	VST1	[V12.S4], (R16)
	ADD	$16, R16, R7
	VST1	[V13.S4], (R7)

	ADD	$8, R13, R13
	ADD	$8, R14, R14
	ADD	$8, R15, R15
	ADD	$32, R16, R16
	ADD	$8, R17, R17
	B	xloop
xdone:
	ADD	R4, R10, R10
	ADD	R5, R2, R2
	ADD	R5, R3, R3
	ADD	R20, R12, R12
	ADD	$1, R8, R8
	B	yloop
done:
	RET

// asmRGBAtoGray converts an *image.RGBA to grayscale float32 pixels. RGB
// channels are de-interleaved with VLD4 and reduced with the integer
// luminosity weights (299*r + 587*g + 114*b) / 1000.
//
// The caller guarantees a width that is a multiple of 8.
//
// ABI0 layout:
//   pixels []float32  pixels+0(FP)
//   pix    []uint8    pix+24(FP)
//   stride            stride+48(FP)
//   minX, minY        minX+56(FP), minY+64(FP)
//   width, height     width+72(FP), height+80(FP)
TEXT ·asmRGBAtoGray(SB), NOSPLIT, $0-88
	MOVD	pixels_base+0(FP), R0
	MOVD	pix_base+24(FP), R1
	MOVD	stride+48(FP), R4
	MOVD	minX+56(FP), R5
	MOVD	minY+64(FP), R6
	MOVD	width+72(FP), R19
	MOVD	height+80(FP), R21

	MUL	R6, R4, R10              // minY * stride
	ADD	R5, R10, R10
	LSL	$2, R5, R7
	ADD	R7, R10, R10             // + minX*4
	ADD	R10, R1, R10             // row base
	LSL	$2, R19, R20             // outStride = width*4
	MOVD	R0, R12                  // op
	MOVD	$0, R8                   // y

	MOVD	$rgbaConst<>(SB), R23
	VLD1R	(R23), [V16.S4]
	ADD	$4, R23, R7
	VLD1R	(R7), [V17.S4]
	ADD	$4, R7, R7
	VLD1R	(R7), [V18.S4]
	ADD	$4, R7, R7
	VLD1R	(R7), [V19.S4]

yloop:
	CMP	R21, R8
	BHS	done
	MOVD	R10, R13                 // row ptr
	MOVD	R12, R14                 // op row
	MOVD	$0, R17                  // x
xloop:
	ADD	$8, R17, R7
	CMP	R19, R7
	BHI	xdone

	VLD4	(R13), [V0.B8, V1.B8, V2.B8, V3.B8]
	VUXTL	V0.B8, V4.H8
	VUXTL	V1.B8, V5.H8
	VUXTL	V2.B8, V6.H8
	VUXTL	V4.H4, V7.S4             // R low
	VUXTL2	V4.H8, V8.S4             // R high
	VUXTL	V5.H4, V9.S4             // G low
	VUXTL2	V5.H8, V10.S4            // G high
	VUXTL	V6.H4, V11.S4            // B low
	VUXTL2	V6.H8, V12.S4            // B high

	VMUL	V16.S4, V7.S4, V7.S4             // 299 * R
	VMUL	V17.S4, V9.S4, V9.S4             // 587 * G
	VMUL	V18.S4, V11.S4, V11.S4           // 114 * B
	VADD	V9.S4, V7.S4, V7.S4
	VADD	V11.S4, V7.S4, V7.S4
	VSCVTF	V7.S4, V7.S4
	VFDIV	V19.S4, V7.S4, V7.S4             // / 1000

	VMUL	V16.S4, V8.S4, V8.S4
	VMUL	V17.S4, V10.S4, V10.S4
	VMUL	V18.S4, V12.S4, V12.S4
	VADD	V10.S4, V8.S4, V8.S4
	VADD	V12.S4, V8.S4, V8.S4
	VSCVTF	V8.S4, V8.S4
	VFDIV	V19.S4, V8.S4, V8.S4

	VST1	[V7.S4], (R14)
	ADD	$16, R14, R15
	VST1	[V8.S4], (R15)

	ADD	$32, R13, R13
	ADD	$32, R14, R14
	ADD	$8, R17, R17
	B	xloop
xdone:
	ADD	R4, R10, R10
	ADD	R20, R12, R12
	ADD	$1, R8, R8
	B	yloop
done:
	RET

DATA ycbcrConst<>+0(SB)/4, $128
DATA ycbcrConst<>+4(SB)/4, $0x10101
DATA ycbcrConst<>+8(SB)/4, $91881
DATA ycbcrConst<>+12(SB)/4, $46802
DATA ycbcrConst<>+16(SB)/4, $22554
DATA ycbcrConst<>+20(SB)/4, $116130
DATA ycbcrConst<>+24(SB)/4, $299
DATA ycbcrConst<>+28(SB)/4, $587
DATA ycbcrConst<>+32(SB)/4, $114
DATA ycbcrConst<>+36(SB)/4, $1000.0
GLOBL ycbcrConst<>(SB), RODATA|NOPTR, $40

DATA rgbaConst<>+0(SB)/4, $299
DATA rgbaConst<>+4(SB)/4, $587
DATA rgbaConst<>+8(SB)/4, $114
DATA rgbaConst<>+12(SB)/4, $1000.0
GLOBL rgbaConst<>(SB), RODATA|NOPTR, $16
