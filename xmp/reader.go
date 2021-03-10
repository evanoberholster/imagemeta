package xmp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/evanoberholster/imagemeta/xmp/xmpns"
)

// xmpRootTag starts with "<x:xmpmeta"
var xmpRootTag = [10]byte{'<', 'x', ':', 'x', 'm', 'p', 'm', 'e', 't', 'a'}

// Reader errors
var (
	ErrNoValue      = errors.New("property has no value")
	ErrNegativeRead = errors.New("error negative read")
	ErrBufferFull   = bufio.ErrBufferFull
)

// Reader blocksizes
const (
	maxTagValueSize  = 256
	maxTagHeaderSize = 64
)

type bufReader struct {
	r *bufio.Reader
	a bool
}

// readRootTag reads and returns the xmpRootTag from the bufReader.
// If the xmpRootTag is not found returns the error ErrNoXMP.
func (br *bufReader) readRootTag() (tag Tag, err error) {
	var buf []byte
	discarded := 0
	for {
		if buf, err = br.Peek(18); err != nil {
			if err == io.EOF {
				err = ErrNoXMP
			}
			return
		}
		if len(buf) < 18 {
			err = ErrNoXMP
			return
		}
		for i := 0; i < 8; i++ {
			if buf[i] == xmpRootTag[0] {
				if bytes.Equal(xmpRootTag[:], buf[i:i+10]) {
					_, err = br.r.ReadSlice('>') // Read until end of the StartTag (RootTag)
					tag.t = startTag
					tag.self = xmpns.XMPRootProperty
					fmt.Println("XMP Discarded:", discarded)
					return tag, err
				}
			}
		}
		discarded += 8
		if _, err = br.Discard(8); err != nil {
			return
		}
	}
}

func (br *bufReader) Discard(n int) (discarded int, err error) {
	return br.r.Discard(n)
}

// hasAttribute returns true when the bufReader's next read is
// an attribute.
func (br *bufReader) hasAttribute() bool {
	return br.a
}

func (br *bufReader) Peek(n int) (buf []byte, err error) {
	if buf, err = br.r.Peek(n); err == io.EOF {
		if len(buf) > 4 {
			return buf, nil
		}
		return buf, err
	}
	return
}

// readAttribute reads an attribute from the bufReader and Tag.
func (br *bufReader) readAttribute(tag *Tag) (attr Attribute, err error) {
	var buf []byte
	attr.pt = attrPType
	attr.parent = tag.self

	// Attribute Name
	if buf, err = br.Peek(maxTagHeaderSize); err != nil {
		err = errors.Wrap(err, "Attr")
		return
	}

	var d int
	if attr.self, d, err = parseAttrName(buf); err != nil {
		err = errors.Wrap(ErrNegativeRead, "Attr (name)")
		return
	}
	if _, err = br.Discard(d); err != nil {
		err = errors.Wrap(err, "Attr (discard)")
		return
	}

	// Attribute Value
	attr.val, err = br.readAttrValue(tag)

	return attr, err
}

// readAttrValue reada an Attributes value from the Tag.
func (br *bufReader) readAttrValue(tag *Tag) (buf []byte, err error) {
	d, i := 0, 2
	s := 64
	for {
		if buf, err = br.Peek(s); err != nil {
			err = errors.Wrap(err, "Attr Value")
			return
		}

		if buf[0] == '=' && (buf[1] == '"' || buf[1] == '\'') {
			delim := buf[1]
			if b := bytes.IndexByte(buf[i:], delim); b >= 0 {
				i += b
				d = i + 1
				if buf[i+1] == '>' {
					d++
					br.a = false
				} else if buf[i+1] == '/' && buf[i+2] == '>' {
					d += 2
					tag.t = soloTag
					br.a = false
				}
				if _, err = br.Discard(d); err != nil {
					err = errors.Wrap(err, "Attr Value (discard)")
				}
				return buf[2:i], err
			}
		}
		s += maxTagValueSize
	}
}

// readTagHeader reads an xmp tag's header and returns the tag.
func (br *bufReader) readTagHeader(parent Tag) (tag Tag, err error) {
	tag.pt = tagPType
	tag.parent = parent.self

	s := maxTagHeaderSize
	// Read Tag Header
	var buf []byte
	var i int
	for {
		if buf, err = br.Peek(s); err != nil {
			err = errors.Wrap(err, "Tag Header")
			return
		}

		// Find Start of Tag
		for ; i < len(buf); i++ {
			if buf[i] == '<' {
				if buf[i+1] == '/' {
					tag.t = stopTag
					i += 2
				} else if buf[i+1] == '?' {
					err = io.EOF
					return
				} else {
					tag.t = startTag
					i++
				}
				buf = buf[i:]
				goto end
			}
		}
		// large white spaces in xmp files

		s += maxTagHeaderSize
	}
end:
	var d int
	tag.self, d, err = parseTagName(buf)
	if err != nil {
		err = errors.Wrap(err, "Tag Header (tag name)") // Err finding tag name
		return
	}
	if buf[d] == '>' {
		br.a = false // No Attributes
		d++
	} else if buf[d] == ' ' || buf[d] == '\n' { // Attributes
		br.a = true
	} else if buf[d] == '/' && buf[d+1] == '>' { // SoloTag
		br.a = false // No Attributes
		tag.t = soloTag
		d += 2
	}
	if _, err = br.Discard(d + i); err != nil {
		err = errors.Wrap(err, "Tag Header (discard)")
	}
	return
}

// readTagValue reads the Tag's Value from the bufReader. Returns
// a temporary []byte.
func (br *bufReader) readTagValue() (buf []byte, err error) {
	var i, j int
	s := maxTagValueSize
	for {
		if buf, err = br.Peek(s); err != nil {
			err = errors.Wrap(err, "Tag Value")
			return
		}
		if i == 0 {
			if buf[i] == '>' {
				i++
			} else if buf[i] == '/' && buf[i+1] == '>' {
				i += 2
			}
			// removes white space and new lines prefixes
			for ; i < len(buf); i++ {
				if buf[i] == ' ' || buf[i] == '\n' {
					continue
				}
				break
			}
			j = i
		}
		// Search buffer.
		for ; j < len(buf); j++ {
			if buf[j] == '<' {
				if _, err = br.Discard(j); err != nil {
					err = errors.Wrap(err, "Tag Value (discard)")
					return nil, err
				}
				return buf[i:j], nil
			}
		}
		s += maxTagValueSize
	}
}

func parseAttrName(buf []byte) (xmpns.Property, int, error) {
	var a, b, c int
	for ; a < len(buf); a++ {
		if buf[a] == ' ' || buf[a] == '\n' {
			continue
		}
		break
	}
	for b = a + 1; b < len(buf); b++ {
		if buf[b] == ':' {
			break
		}
	}
	for c = b + 2; c < len(buf); c++ {
		if buf[c] == '=' || buf[c] == ' ' {
			return xmpns.IdentifyProperty(buf[a:b], buf[b+1:c]), c, nil
		}
	}
	return xmpns.Property{}, -1, ErrNegativeRead
}

func parseTagName(buf []byte) (xmpns.Property, int, error) {
	var a, b int
	for ; a < len(buf); a++ {
		if buf[a] == ':' {
			break
		}
	}
	for b = a + 1; b < len(buf); b++ {
		if buf[b] == '>' || buf[b] == ' ' || buf[b] == '\n' || buf[b] == '/' {
			return xmpns.IdentifyProperty(buf[:a], buf[a+1:b]), b, nil
		}
	}
	return xmpns.Property{}, -1, ErrNegativeRead
}
