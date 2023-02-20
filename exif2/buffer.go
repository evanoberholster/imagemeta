package exif2

import (
	"sync"

	"github.com/rs/zerolog"
)

const (
	tagMaxCount  = 84
	bufferLength = 1024
)

// buffer for data and tags
type buffer struct {
	buf [bufferLength]byte
	tag [tagMaxCount]Tag
	len uint32
	pos uint32
}

// bufferPool for tag buffers
var bufferPool = sync.Pool{
	New: func() interface{} { return new(buffer) },
}

// currentTag returns the current tag in tagBuffer
func (b *buffer) currentTag() Tag {
	return b.tag[b.pos]
}

// nextTag returns the next tag in tagBuffer
func (b *buffer) nextTag() Tag {
	return b.tag[b.pos+1]
}

// nextTag increments the position by 1
func (b *buffer) advanceBuffer() Tag {
	if b.pos < b.len {
		b.pos++
		return b.tag[b.pos]
	}
	return Tag{}
}

// validTag returns true if the tag is valid
func (b *buffer) validTag() bool {
	return b.pos < b.len && b.len > 0
}

// readTagValue discards until tag.ValueOffset and reads length of tag
func (ir *ifdReader) readTagValue() (buf []byte, err error) {
	t := ir.buffer.currentTag()
	if err := ir.discard(int(t.ValueOffset) - int(ir.po)); err != nil {
		return nil, err
	}
	return ir.fastRead(int(t.Size()))
}

// seekToTag seeks with the underlying reader to given tag value
func (ir *ifdReader) seekToTag(t Tag) (err error) {
	discard := int(t.ValueOffset) - int(ir.po)
	if err = ir.discard(discard); err != nil && ir.logLevelError() {
		t.logTag(ir.logError(err)).Uint32("ifdReaderPosition", ir.po).Uint32("discard", uint32(discard)).Send()
	}
	return
}

// resetPosition resets the tag buffer to only include unread tags
func (b *buffer) resetPosition() {
	if b.pos > 0 {
		copy(b.tag[:b.len-b.pos], b.tag[b.pos:b.len])
		b.len -= b.pos
		b.pos = 0
	}
}

// addTagBuffer adds the given tag to the tagBuffer
func (ir *ifdReader) addTagBuffer(t Tag) {
	if t.ValueOffset < ir.po {
		if ir.logLevelWarn() {
			t.logTag(ir.logWarn()).Uint32("readerOffset", ir.po).Msg("Incompatible reverse exif tag")
		}
		return
	}
	b := ir.buffer
	if b.len < tagMaxCount {
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
	if ir.logLevelWarn() {
		ir.logWarn().Int32("tagMaxCount", tagMaxCount).Msg("error tagMaxCount is too short")
	}
}

// discard, discards n amount from ir.Reader
func (ir *ifdReader) discard(n int) (err error) {
	if n == 0 {
		return nil
	}
	if int(ir.exifLength) < n+int(ir.po) {
		n = int(ir.exifLength) - int(ir.po)
	}
	if ir.bufReader != nil {
		n, err := ir.bufReader.Discard(n)
		ir.po += uint32(n)
		return err
	}
	var discarded int
	for n > 0 && err == nil {
		if bufferLength > n {
			discarded, err = ir.reader.Read(ir.buffer.buf[:n])
		} else {
			discarded, err = ir.reader.Read(ir.buffer.buf[:])
		}
		ir.po += uint32(discarded)
		n -= discarded
	}
	return err
}

// clear the buffer counters
func (b *buffer) clear() {
	b.len = 0
	b.pos = 0
}

// MarshalZerologArray is a zerolog interface for logging
func (b *buffer) MarshalZerologArray(a *zerolog.Array) {
	for i := b.pos; i < b.len; i++ {
		a.Object(b.tag[i])
	}
}
