package tag

import (
	"errors"
)

// Errors
var (
	ErrEmptyTag      = errors.New("Error empty tag")
	ErrTagNotValid   = errors.New("Error tag not valid")
	ErrNotEnoughData = errors.New("Error not enough data to parse tag")
)

// RawEncodedBytes returns the raw encoded bytes for the value that we represent.
func (t Tag) RawEncodedBytes(tagReader TagReader) (rawBytes []byte, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	//tagType := tag.effectiveValueType()

	byteLength := t.TagType.Size() * t.UnitCount

	// check if Value is Embedded
	if valueIsEmbbeded(byteLength) {
		return t.rawValueOffset[:byteLength], nil
	}

	data := make([]byte, byteLength)
	if _, err = tagReader.ReadAt(data, int64(t.valueOffset)); err != nil {
		panic(err)
	}
	return data, nil
}

// ASCIIValue returns the ASCII value of the tag as a string
// and returns an error if it encounters one
func (t *Tag) ASCIIValue(tagReader TagReader) (value string, err error) {
	if t.TagType.IsValid() {
		// Needs Typecheck
		var rawBytes []byte
		rawBytes, err = t.RawEncodedBytes(tagReader)
		if err != nil {
			return
		}

		// Trim trailing spaces 0x00 and 0x20
		rawBytes = trim(rawBytes)
		value = string(rawBytes)

		return
	}
	err = ErrTagNotValid
	return
}

// Uint16Value returns the Short value of the tag as a uint16
// and returns an error if it encounters one.
//
// Warning: it returns only the first value if there are more values
// use Uint16Values function
func (t *Tag) Uint16Value(tagReader TagReader) (value uint16, err error) {
	if t.TagType == TypeShort {
		var rawBytes []byte
		rawBytes, err = t.RawEncodedBytes(tagReader)
		if err != nil {
			return 0, err
		}
		if len(rawBytes) < 2 {
			err = ErrEmptyTag
			return
		}
		value = tagReader.ByteOrder().Uint16(rawBytes[:2])

		return
	}
	err = ErrTagTypeNotValid
	return
}

// Uint16Values returns the Short value of the tag as a uint16 array
// and returns an error if it encounters one.
func (t *Tag) Uint16Values(tagReader TagReader) (value []uint16, err error) {
	if t.TagType == TypeShort {

		var rawBytes []byte
		if rawBytes, err = t.RawEncodedBytes(tagReader); err != nil {
			return nil, err
		}

		byteOrder := tagReader.ByteOrder()
		count := int(t.UnitCount)

		if len(rawBytes) < (TypeShortSize * count) {
			err = ErrNotEnoughData
		}

		value = make([]uint16, count)
		for i := 0; i < count; i++ {
			value[i] = byteOrder.Uint16(rawBytes[i*2:])
		}

		return
	}
	err = ErrTagTypeNotValid
	return
}

// Uint32Values returns the Long value of the tag as a uint32 array
// and returns an error if it encounters one.
//
func (t *Tag) Uint32Values(tagReader TagReader) (value []uint32, err error) {
	if t.TagType == TypeLong {

		var rawBytes []byte
		if rawBytes, err = t.RawEncodedBytes(tagReader); err != nil {
			return nil, err
		}

		byteOrder := tagReader.ByteOrder()
		count := int(t.UnitCount)

		if len(rawBytes) < (TypeLongSize * count) {
			err = ErrNotEnoughData
		}

		value = make([]uint32, count)
		for i := 0; i < count; i++ {
			value[i] = byteOrder.Uint32(rawBytes[i*4:])
		}

		return
	}
	err = ErrTagTypeNotValid
	return
}

// RationalValues returns a list of unsignedRationals
func (t *Tag) RationalValues(tagReader TagReader) (value []Rational, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()
	if t.TagType == TypeRational || t.TagType == TypeSignedRational {
		var rawBytes []byte
		if rawBytes, err = t.RawEncodedBytes(tagReader); err != nil {
			return nil, err
		}

		byteOrder := tagReader.ByteOrder()
		count := int(t.UnitCount)

		if len(rawBytes) < (TypeRationalSize * count) {
			panic(ErrNotEnoughData)
		}

		value = make([]Rational, count)
		for i := 0; i < count; i++ {
			value[i].Numerator = byteOrder.Uint32(rawBytes[i*8:])
			value[i].Denominator = byteOrder.Uint32(rawBytes[i*8+4:])
		}

		return
	}
	err = ErrTagTypeNotValid
	return
}

////
// Helper functions
////

// valueIsEmbbeded checks if Tag value is embedded in the Tag.RawValueOffset
func valueIsEmbbeded(byteLength uint32) bool {
	return byteLength <= 4
}

// trim removes trailing 0x00 and 0x20 from []byte
func trim(buf []byte) []byte {
	i := len(buf)
	for i > 0 {
		if buf[i-1] == 0x0 ||
			buf[i-1] == 0x20 {
			i--
		} else {
			break
		}
	}
	return buf[:i]
}
