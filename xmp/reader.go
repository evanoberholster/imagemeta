package xmp

import (
	"bufio"
	"bytes"
	"io"

	"github.com/pkg/errors"

	"github.com/evanoberholster/imagemeta/xmp/xmpns"
)

// xmpRootTag starts with "<x:xmpmeta"
var xmpRootTag = [10]byte{'<', 'x', ':', 'x', 'm', 'p', 'm', 'e', 't', 'a'}

type bufReader struct {
	r *bufio.Reader
	a bool
}

// readRootTag reads and returns the xmpRootTag from the bufReader.
// If the xmpRootTag is not found returns the error ErrNoXMP.
func (br *bufReader) readRootTag() (tag Tag, err error) {
	var buf []byte
	for {
		if buf, err = br.Peek(18); err != nil {
			if err == io.EOF {
				err = ErrNoXMP
			}
			return
		}
		for i := 0; i < 8; i++ {
			if buf[i] == xmpRootTag[0] {
				if bytes.Equal(xmpRootTag[:], buf[i:i+10]) {
					_, err = br.r.ReadSlice('>') // Read until end of the StartTag (RootTag)
					tag.t = startTag
					tag.self = xmpns.XMPRootProperty
					return tag, err
				}
			}
		}
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
	if buf, err = br.Peek(maxTagHeaderSize); err != nil {
		return
	}
	var i int
	for i = 0; i < len(buf); i++ {
		if buf[i] == ' ' || buf[i] == '\n' {
			continue
		}
		if buf[i] == '>' {
			br.a = false
			i++
			return
		} else if buf[i] == '/' && buf[i+1] == '>' {
			br.a = false
			i += 2
			tag.t = soloTag
			return
		}
		break
	}
	buf = buf[i:]

	// Attribute Name
	var a, b int
	if a, b = attrNameIndex(buf); a == -1 {
		err = ErrNegativeRead
		return
	}
	attr.self = xmpns.IdentifyProperty(buf[:a], buf[a+1:b])

	if _, err = br.r.Discard(i + b); err != nil {
		return
	}

	attr.val, err = br.readAttrValue(tag)
	return attr, err
}

// readAttrValue reada an Attributes value from the Tag.
func (br *bufReader) readAttrValue(tag *Tag) (buf []byte, err error) {
	// Attribute Value
	s := maxTagValueSize
	var i int
	for {
		if buf, err = br.Peek(s); err != nil {
			return
		}

		var delim byte
		delim = '"'
		if buf[0] == '=' {
			delim = buf[1]
			buf = buf[2:]
		}
		for ; i < len(buf); i++ {
			if buf[i] == delim {
				goto end
			}
		}
		s += maxTagValueSize
	}
end:
	if buf[i+1] == '>' {
		br.a = false
	} else if buf[i+1] == '/' && buf[i+2] == '>' {
		tag.t = soloTag
		br.a = false
	}
	if _, err = br.Discard(i + 3); err != nil {
		return
	}
	return buf[:i], nil
}

// readTagHeader reads an xmp tag's header and returns the tag.
func (br *bufReader) readTagHeader(parent Tag) (tag Tag, err error) {
	tag.pt = tagPType
	tag.parent = parent.self

	// Read Tag Header
	var buf []byte
	var i int
	for {
		if buf, err = br.Peek(maxTagHeaderSize); err != nil {
			err = errors.Wrap(err, "Tag Header")
			return
		}

		// Find Start of Tag
		if i = bytes.IndexByte(buf, '<'); i >= 0 {
			break
		}

		if _, err = br.Discard(maxTagHeaderSize); err != nil {
			err = errors.Wrap(err, "Tag Header (discard)")
			return
		}
	}

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
	buf = buf[i:] // reslice tag
	var a, b int
	if a, b = tagNameIndex(buf); a < 0 {
		err = errors.Wrap(ErrNegativeRead, "Tag Header (tag name)") // Err finding tag name
		return
	}
	tag.self = xmpns.IdentifyProperty(buf[:a], buf[a+1:b])
	if buf[b] == '>' {
		br.a = false // No Attributes
		b++
	} else if buf[b] == ' ' || buf[b] == '\n' { // Attributes
		br.a = true
	} else if buf[b] == '/' && buf[b+1] == '>' { // SoloTag
		br.a = false // No Attributes
		tag.t = soloTag
		b += 2
	}
	if _, err = br.Discard(b + i); err != nil {
		err = errors.Wrap(err, "Tag Header (discard)")
		return // error here
	}
	return
}

// readTagValue reads the Tag's Value from the bufReader. Returns
// a temporary []byte.
func (br *bufReader) readTagValue() (buf []byte, err error) {
	var i int
	s := maxTagValueSize
	for {
		if buf, err = br.Peek(s); err != nil {
			return nil, err
		}
		if buf[0] == '>' {
			buf = buf[1:]
		} else if buf[0] == '/' && buf[1] == '>' {
			buf = buf[2:]
		}
		// Search buffer.
		if i = bytes.IndexByte(buf, '<'); i >= 0 {
			buf = buf[:i]
			break
		}
		s += maxTagValueSize
	}
	if _, err = br.Discard(i); err != nil {
		return nil, err
	}

	// Remove white space and new lines prefix
	for i = 0; i < len(buf); i++ {
		if buf[i] == ' ' || buf[i] == '\n' {
			continue
		}
		break
	}
	return buf[i:], nil
}

func attrNameIndex(buf []byte) (int, int) {
	var a, i int
	for ; i < len(buf); i++ {
		if buf[i] == ':' {
			a = i
			break
		}
	}
	for ; i < len(buf); i++ {
		if buf[i] == '=' || buf[i] == ' ' {
			return a, i
		}
	}
	return -1, -1
}

func tagNameIndex(buf []byte) (int, int) {
	var a, i int
	for ; i < len(buf); i++ {
		if buf[i] == ':' {
			a = i
			break
		}
	}
	for ; i < len(buf); i++ {
		if buf[i] == '>' || buf[i] == ' ' || buf[i] == '\n' || buf[i] == '/' {
			return a, i
		}
	}
	return -1, -1
}
