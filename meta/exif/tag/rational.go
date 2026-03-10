package tag

// RationalU stores an unsigned rational number.
type RationalU struct {
	Numerator   uint32
	Denominator uint32
}

// Float64 converts the rational value into a float64.
func (r RationalU) Float64() float64 {
	if r.Denominator == 0 {
		return 0
	}
	return float64(r.Numerator) / float64(r.Denominator)
}
