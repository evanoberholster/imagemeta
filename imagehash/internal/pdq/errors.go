package pdq

import "errors"

var (
	// ErrNilImage is returned when a nil image is passed to Hash.
	ErrNilImage = errors.New("pdq: image must not be nil")

	ErrMatrixLength        = errors.New("pdq: matrix length is not the expected 64x64")
	ErrTorbenElementLength = errors.New("pdq: expected 256 elements for torben median")
	ErrQuantizeLength      = errors.New("pdq: expected 16x16 for DCT quantization")
	ErrDihedralDCTSize     = errors.New("pdq: expected 16x16 for DCT dihedral")
)
