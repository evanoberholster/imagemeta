package exif2

import (
	"bufio"
	"fmt"
	"strings"
	"sync"

	"github.com/evanoberholster/imagemeta/exif2/ifds"
	"github.com/evanoberholster/imagemeta/exif2/tag"
)

const (
	tagMaxLength = 56
	bufferLength = 512
)

// buffer for data and tags
type buffer struct {
	buf [bufferLength]byte
	tag [tagMaxLength]tag.Tag
	len uint32
	pos uint32
}

// bufferPool for buffer
var bufferPool = sync.Pool{
	New: func() interface{} { return new(buffer) },
}

// currentTag returns the current tag in tagBuffer
func (b *buffer) currentTag() tag.Tag {
	return b.tag[b.pos]
}

// nextTag returns the next tag in tagBuffer
func (b *buffer) nextTag() tag.Tag {
	return b.tag[b.pos+1]
}

// nextTag increments the position by 1
func (b *buffer) advanceBuffer() tag.Tag {
	if b.pos < b.len {
		b.pos++
		return b.tag[b.pos]
	}
	return tag.Tag{}
}

// validTag returns true if the tag is valid
func (b *buffer) validTag() bool {
	return b.pos < b.len && b.len > 0
}

// readTagValue discards until tag.ValueOffset and reads length of tag
func (ir *ifdReader) readTagValue() ([]byte, error) {
	t := ir.buffer.currentTag()
	if err := ir.discard(int(t.ValueOffset) - int(ir.po)); err != nil {
		return nil, err
	}
	//
	br, ok := ir.reader.(*bufio.Reader)
	if ok {
		buf, err := br.Peek(int(t.Size()))
		if err != nil {
			panic(err)
		}
		n, err := br.Discard(len(buf))
		ir.po += uint32(n)
		return buf, err
	}
	n, err := ir.reader.Read(ir.buffer.buf[:t.Size()])
	ir.po += uint32(n)
	return ir.buffer.buf[:n], err
}

// resetPosition resets the tag buffer to only include unread tags
func (b *buffer) resetPosition() {
	if b.pos > 0 {
		copy(b.tag[:b.len-b.pos], b.tag[b.pos:b.len])
		b.len -= b.pos
		b.pos = 0
		//fmt.Println("Position: ", ir.po, ir.tagBuf.pos, ir.tagBuf.len)
	}
}

// addTagBuffer adds the given tag to the tagBuffer
func (ir *ifdReader) addTagBuffer(t tag.Tag) {
	if uint32(t.ValueOffset) < ir.po {
		ir.logTagWarn(t, "Uncompatible reverse exif tag")
		return
	}
	b := ir.buffer
	if b.len < tagMaxLength {
		for i := b.len; i > 0; i-- {
			if t.ValueOffset > b.tag[i-1].ValueOffset {
				if i != b.len {
					copy(b.tag[i+1:b.len+1], b.tag[i:b.len])
				}
				b.tag[i] = t
				b.len++
				return
			}
		}
		if b.len == 0 {
			b.tag[0] = t
			b.len++
			return
		}
		if t.ValueOffset < b.tag[0].ValueOffset {
			copy(b.tag[0+1:b.len+1], b.tag[0:b.len])
			b.tag[0] = t
			b.len++
			return
		}
	}
	if ir.logError() {
		ir.logger.Error().Int32("tagBufferLength", tagMaxLength).Msg("error tagBufferMaxLength is too short")
	}
}

// discard, discards n amount from ir.Reader
func (ir *ifdReader) discard(n int) error {
	if n == 0 {
		return nil
	}
	if int(ir.exifLength) < n+int(ir.po) {
		n = int(ir.exifLength) - int(ir.po)
	}
	var discarded int
	var err error
	for n > 0 {
		if bufferLength > n {
			discarded, err = ir.reader.Read(ir.buffer.buf[:n])
		} else {
			discarded, err = ir.reader.Read(ir.buffer.buf[:])
		}
		ir.po += uint32(discarded)
		n -= discarded
		if err != nil {
			return err
		}
	}
	return err
}

// clear the buffer counters
func (b *buffer) clear() {
	b.len = 0
	b.pos = 0
}

// discard, discards the given amount from ir.Reader
//func (ir *ifdReader) discard2(n int) error {
//	if n == 0 {
//		return nil
//	}
//
//	discarded, err := io.CopyN(io.Discard, ir.reader, int64(n))
//	ir.po += uint32(discarded)
//	return err
//}

// String is the stringer interface for buffer
func (b *buffer) String() string {
	sb := strings.Builder{}
	for i := b.pos; i < b.len; i++ {
		sb.WriteString(tagString(b.tag[i]) + "\n")
	}
	return sb.String()
}

func tagString(t tag.Tag) string {
	return fmt.Sprintf("%s\t | %s \t\t| %s \t|%d", t.ID, ifds.IfdType(t.Ifd).TagName(t.ID), t.Type(), t.ValueOffset)
}
