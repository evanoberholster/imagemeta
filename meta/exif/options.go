package exif

// AFInfoDecodeOptions controls which Canon AFInfo2/AFInfo3 structures are
// materialized during parse.
type AFInfoDecodeOptions uint8

const (
	// AFInfoDecodeCoords decodes AFArea width/height/x/y coordinate lists.
	AFInfoDecodeCoords AFInfoDecodeOptions = 1 << iota
	// AFInfoDecodeInFocus decodes AFPointsInFocusBits/InFocus.
	AFInfoDecodeInFocus
	// AFInfoDecodeSelected decodes AFPointsSelectedBits/Selected for EOS models.
	AFInfoDecodeSelected
	// AFInfoDecodePoints decodes derived AFPoints rectangles.
	AFInfoDecodePoints
)

const (
	// AFInfoDecodeAll preserves the historical/full decode behavior.
	AFInfoDecodeAll = AFInfoDecodeCoords | AFInfoDecodeInFocus | AFInfoDecodeSelected | AFInfoDecodePoints
)

func (o AFInfoDecodeOptions) has(flag AFInfoDecodeOptions) bool {
	return o&flag != 0
}

// ReaderOption configures Reader behavior.
type ReaderOption func(*Reader)

func newSetAFInfoDecodeOption(opts AFInfoDecodeOptions) ReaderOption {
	return func(r *Reader) {
		r.afInfoDecodeOptions = opts
	}
}

var afInfoDecodeOptionTable = func() [16]ReaderOption {
	var table [16]ReaderOption
	for i := range table {
		table[i] = newSetAFInfoDecodeOption(AFInfoDecodeOptions(i))
	}
	return table
}()

// WithAFInfoDecodeOptions sets Canon AFInfo decode flags for this reader.
func WithAFInfoDecodeOptions(opts AFInfoDecodeOptions) ReaderOption {
	// Keep options in the supported bitset range and return a precomputed
	// function to avoid allocating a fresh closure on each call.
	return afInfoDecodeOptionTable[int(opts&AFInfoDecodeAll)]
}

func applyReaderOptions(r *Reader, opts []ReaderOption) {
	switch len(opts) {
	case 0:
		return
	case 1:
		if opt := opts[0]; opt != nil {
			opt(r)
		}
		return
	}
	for i := 0; i < len(opts); i++ {
		if opt := opts[i]; opt != nil {
			opt(r)
		}
	}
}
