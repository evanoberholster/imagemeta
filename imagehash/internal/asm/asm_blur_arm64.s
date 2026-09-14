//go:build arm64

#include "textflag.h"

// asmBlurRow computes the per-row BlurHash basis sums for a single 64-wide
// image row. Each NEON lane accumulates one x-component (xc), so the four
// components for all three channels are produced without a horizontal
// reduction.
//
// out is laid out as out[channel*xComponents+xc].
//
// ABI0 layout:
//   out      []float32  out+0(FP)
//   lr       []float32  lr+24(FP)
//   lg       []float32  lg+48(FP)
//   lb       []float32  lb+72(FP)
//   xvaluesT []float32  xvaluesT+96(FP)
TEXT ·asmBlurRow(SB), NOSPLIT, $0-120
	MOVD	out_base+0(FP), R0
	MOVD	lr_base+24(FP), R1
	MOVD	lg_base+48(FP), R2
	MOVD	lb_base+72(FP), R3
	MOVD	xvaluesT_base+96(FP), R4

	VEOR	V0.B16, V0.B16, V0.B16
	VEOR	V1.B16, V1.B16, V1.B16
	VEOR	V2.B16, V2.B16, V2.B16
	MOVD	$64, R5
loop:
	VLD1R	(R1), [V3.S4]            // lr[x] broadcast
	VLD1R	(R2), [V4.S4]            // lg[x]
	VLD1R	(R3), [V5.S4]            // lb[x]
	VLD1	(R4), [V6.S4]            // xvaluesT[x][0..3]
	VFMLA	V3.S4, V6.S4, V0.S4      // accR += lr*xv
	VFMLA	V4.S4, V6.S4, V1.S4      // accG += lg*xv
	VFMLA	V5.S4, V6.S4, V2.S4      // accB += lb*xv
	ADD	$4, R1, R1
	ADD	$4, R2, R2
	ADD	$4, R3, R3
	ADD	$16, R4, R4
	SUB	$1, R5, R5
	CBNZ	R5, loop

	VST1	[V0.S4], (R0)
	ADD	$16, R0, R6
	VST1	[V1.S4], (R6)
	ADD	$16, R6, R6
	VST1	[V2.S4], (R6)
	RET
