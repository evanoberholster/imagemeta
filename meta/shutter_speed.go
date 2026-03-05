package meta

import "strconv"

// ShutterSpeed stores shutter duration in seconds.
type ShutterSpeed float32

func (ss ShutterSpeed) String() string {
	return strconv.FormatFloat(float64(ss), 'f', -1, 32)
}
