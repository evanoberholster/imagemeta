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

// Package bmff reads ISO BMFF boxes, as used by HEIF, AVIF, CR3, etc.
package bmff

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

// Debug Flag
var (
	Debug = false
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
func (r *Reader) ReadMetaBox() (MetaBox, error) {
	if r.brand == brandUnknown {
		return MetaBox{}, fmt.Errorf("brand not supported")
	}
	if r.noMoreBoxes {
		return MetaBox{}, fmt.Errorf("no more boxes to be parsed")
	}
	b, err := r.br.readInnerBox()
	if err != nil {
		panic(err)
	}
	return parseMetaBox(&b)
}

// ReadBox reads a box and returns it
func (r *Reader) readBox() (b box, err error) {
	var buf []byte
	// Read box size and box type
	if buf, err = r.br.Peek(8); err != nil {
		return b, err
	}

	b = box{
		bufReader: r.br,
		size:      int64(binary.BigEndian.Uint32(buf[:4])),
		boxType:   boxType(buf[4:8]),
	}

	if err = r.br.discard(8); err != nil {
		return
	}

	var remain int
	switch b.size {
	case 1:
		// 1 means it's actually a 64-bit size, after the type.
		if buf, err = r.br.Peek(8); err != nil {
			return b, err
		}
		b.size = int64(binary.BigEndian.Uint64(buf[:8]))
		if b.size < 0 {
			// Go uses int64 for sizes typically, but BMFF uses uint64.
			// We assume for now that nobody actually uses boxes larger
			// than int64.
			return b, fmt.Errorf("unexpectedly large box %q", b.boxType)
		}
		remain = int(b.size - 2*4 - 8)
		if err = r.br.discard(8); err != nil {
			return
		}
	case 0:
		// 0 means unknown & to read to end of file. No more boxes.
		//r.noMoreBoxes = true
	default:
		remain = int(b.size - 2*4)
	}
	b.remain = remain
	return b, nil
}

// Box represents a BMFF box.
type Box interface {
	//Size() int64 // 0 means unknown (will read to end of file)
	Type() BoxType

	// Parses parses the box, populating the fields
	// in the returned concrete type.
	//
	// If Parse has already been called, Parse returns nil.
	// If the box type is unknown, the returned error is ErrUnknownBox
	// and it's guaranteed that no bytes have been read from the box.
	//Parse() (Box, error)

	// Body returns the inner bytes of the box, ignoring the header.
	// The body may start with the 4 byte header of a "Full Box" if the
	// box's type derives from a full box. Most users will use Parse
	// instead.
	// Body will return a new reader at the beginning of the box if the
	// outer box has already been parsed.
	//Body() io.Reader
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