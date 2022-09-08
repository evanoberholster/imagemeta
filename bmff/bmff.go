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

// Package bmff reads ISOBMFF boxes, as used by HEIF, AVIF, CR3, etc. and other riff based files
package bmff

import (
	"bufio"
	"io"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/pkg/errors"
)

// Common Errors
var (
	ErrBrandNotSupported = errors.New("error brand not supported")
	ErrWrongBoxType      = errors.New("error wrong box type")
	ErrNoMoreBoxes       = errors.New("no more boxes to be parsed")
)

// brandCount is the number of compatible brands supported.
const brandCount = 8

// NewReader returns a new bmff.Reader
func NewReader(r io.Reader) Reader {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	return Reader{br: bufReader{Reader: br, remain: 8}}
}

// Reader is a BMFF reader
type Reader struct {
	br          bufReader
	brand       Brand
	noMoreBoxes bool // a box with size 0 (the final box) was seen
	ExifReader  func(r io.Reader, h meta.ExifHeader) error
}

// ReadAndParseBox wraps the ReadBox method, ensuring that the read box is of type typ
// and parses successfully. It returns the parsed box.
func (r *Reader) ReadAndParseBox(typ BoxType) (Box, error) {
	box, err := r.readBox()
	if err != nil {
		return nil, errors.Errorf("error reading %q box: %v", typ, err)
	}
	if box.Type() != typ {
		return nil, errors.Errorf("error reading %q box: got box type %q instead", typ, box.Type())
	}
	pbox, err := box.Parse()
	if err != nil {
		return nil, errors.Errorf("error parsing read %q box: %v", typ, err)
	}
	return pbox, nil
}

// ReadFtypBox reads an 'ftyp' box from a BMFF file.
//
// This should be the first read function called.
func (r *Reader) ReadFtypBox() (FileTypeBox, error) {
	b, err := r.readBox()
	if err != nil {
		return FileTypeBox{}, errors.Wrapf(err, "ReadFtypBox")
	}
	ftyp, err := b.parseFileTypeBox()
	r.brand = ftyp.MajorBrand
	r.br.offset = b.offset
	return ftyp, err
}

// ReadMetaBox reads a 'meta' box from a BMFF file.
//
// This should be called in order. First call ReadFtypBox
func (r *Reader) ReadMetaBox() (mb MetaBox, err error) {
	if r.brand == brandUnknown {
		return mb, ErrBrandNotSupported
	}
	if r.noMoreBoxes {
		return mb, ErrNoMoreBoxes
	}
	b, err := r.readBox()
	if err != nil {
		err = errors.Wrapf(err, "ReadMetaBox")
		return mb, err
	}
	return parseMetaBox(&b)
}

// ReadMoovBox reads a 'moov' box from a BMFF file.
//
// This should be called in order. First call ReadFtypBox
func (r *Reader) ReadMoovBox() (moov MoovBox, err error) {
	if r.brand == brandUnknown {
		return moov, ErrBrandNotSupported
	}
	if r.noMoreBoxes {
		return moov, ErrNoMoreBoxes
	}
	b, err := r.readBox()
	if err != nil {
		err = errors.Wrapf(err, "ReadMetaBox")
		return moov, err
	}
	return b.parseMoovBox()
}

// ReadBox reads a box and returns it
func (r *Reader) readBox() (b box, err error) {
	return r.br.readInnerBox()
}
