package canon

import "fmt"

var (
	_ fmt.Stringer = ContinuousDrive(0)
	_ fmt.Stringer = FocusMode(0)
	_ fmt.Stringer = MeteringMode(0)
	_ fmt.Stringer = FocusRange(0)
	_ fmt.Stringer = ExposureMode(0)
	_ fmt.Stringer = BracketMode(0)
	_ fmt.Stringer = AESetting(0)
	_ fmt.Stringer = AFAreaMode(0)
)

