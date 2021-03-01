/*
	Copyright 2021 Evan Oberholster
	Copyright 2018 The go4 Authors

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

// Package bmff reads ISO BMFF boxes, as used by HEIF, AVIF, CR3, etc. and other riff based files
package bmff

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Debug Flag
var (
	Debug = false
)

// Common Errors
var (
	ErrBrandNotSupported = errors.New("error brand not supported")
	ErrWrongBoxType      = errors.New("error wrong box type")
)

// NewReader returns a new bmff.Reader
func NewReader(r io.Reader) Reader {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	return Reader{br: bufReader{Reader: br}}
}

// Reader is a BMFF reader
type Reader struct {
	br          bufReader
	brand       Brand
	noMoreBoxes bool // a box with size 0 (the final box) was seen
}

// ReadAndParseBox wraps the ReadBox method, ensuring that the read box is of type typ
// and parses successfully. It returns the parsed box.
func (r *Reader) ReadAndParseBox(typ BoxType) (Box, error) {
	box, err := r.readBox()
	if err != nil {
		return nil, fmt.Errorf("error reading %q box: %v", typ, err)
	}
	if box.Type() != typ {
		return nil, fmt.Errorf("error reading %q box: got box type %q instead", typ, box.Type())
	}
	pbox, err := box.Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing read %q box: %v", typ, err)
	}
	return pbox, nil
}

// ReadFtypBox reads an 'ftyp' box from a BMFF file.
//
// This should be the first read function called.
func (r *Reader) ReadFtypBox() (FileTypeBox, error) {
	b, err := r.br.readInnerBox()
	if err != nil {
		return FileTypeBox{}, err
	}
	ftyp, err := parseFileTypeBox(&b)
	r.brand = ftyp.MajorBrand
	return ftyp, err
}

// ReadMetaBox reads a 'meta' box from a BMFF file.
//
// This should be called in order. First call ReadFtypBox
func (r *Reader) ReadMetaBox() (mb MetaBox, err error) {
	if r.brand == brandUnknown {
		err = ErrBrandNotSupported
		return
	}
	if r.noMoreBoxes {
		err = fmt.Errorf("no more boxes to be parsed")
		return
	}
	b, err := r.br.readInnerBox()
	if err != nil {
		return mb, err
	}
	return parseMetaBox(&b)
}

// ReadMoovBox reads a 'moov' box from a BMFF file.
//
// This should be called in order. First call ReadFtypBox
func (r *Reader) ReadMoovBox() (moov MoovBox, err error) {
	if r.brand == brandUnknown {
		err = ErrBrandNotSupported
		return
	}
	if r.noMoreBoxes {
		err = fmt.Errorf("no more boxes to be parsed")
		return
	}
	b, err := r.br.readInnerBox()
	if err != nil {
		return moov, err
	}
	return parseMoovBox(&b)
}

// ReadBox reads a box and returns it
func (r *Reader) readBox() (b box, err error) {
	outer := box{bufReader: r.br}
	return outer.readInnerBox()
}

// Box represents a BMFF box.
type Box interface {
	//Size() int64 // 0 means unknown (will read to end of file)
	Type() BoxType
}

// Flags for a FullBox
// 8 bits -> Version
// 24 bits -> Flags
type Flags uint32

// Flags returns underlying Flags after removing version.
// Flags are 24 bits.
func (f Flags) Flags() uint32 {
	// Left Shift
	f = f << 8
	// Right Shift
	return uint32(f >> 8)
}

// Version returns a uint8 version.
func (f Flags) Version() uint8 {
	return uint8(f >> 24)
}
