package exif

import (
	"sync"

	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

const (
	readBufferLength = 4096
	tagQueueMax      = 128
)

type state struct {
	buf        [readBufferLength]byte
	discardBuf [readBufferLength]byte
	tag        [tagQueueMax]tag.Entry
	dirty      bool
	len        uint32
	pos        uint32
}

var statePool = sync.Pool{
	New: func() any { return new(state) },
}

// reset resets parser state for a new decode operation.
func (s *state) reset() {
	s.len = 0
	s.pos = 0
	s.dirty = false
}

// currentTag returns the current tag entry from the parser queue.
func (s *state) currentTag() tag.Entry {
	return s.tag[s.pos]
}

// advanceTag advances and returns the next tag entry from the parser queue.
func (s *state) advanceTag() tag.Entry {
	if s.pos+1 < s.len {
		s.pos++
		return s.tag[s.pos]
	}
	s.pos = s.len
	return tag.Entry{}
}

// validTag reports whether the current queued tag is valid for parsing.
func (s *state) validTag() bool {
	return s.pos < s.len && s.len > 0
}

// resetPosition compacts unread tags at the front of the queue.
func (s *state) resetPosition() {
	if s.pos == 0 {
		return
	}
	copy(s.tag[:s.len-s.pos], s.tag[s.pos:s.len])
	s.len -= s.pos
	s.pos = 0
}

// addTag appends a tag to parser state and marks queue sort state as needed.
func (s *state) addTag(t tag.Entry) bool {
	if s.len >= tagQueueMax {
		return false
	}

	if s.len == 0 {
		s.tag[0] = t
		s.len = 1
		s.dirty = false
		return true
	}

	last := s.tag[s.len-1].ValueOffset
	s.tag[s.len] = t
	s.len++
	if !s.dirty && t.ValueOffset < last {
		s.dirty = true
	}
	return true
}

// sortAll sorts all queued tags by ValueOffset.
func (s *state) sortAll() {
	if !s.dirty || s.len < 2 {
		return
	}
	s.sortRange(0, s.len)
	s.dirty = false
}

// sortUnread sorts unread queued tags while preserving the current tag position.
func (s *state) sortUnread() {
	if !s.dirty {
		return
	}
	start := s.pos + 1
	if start >= s.len {
		s.dirty = false
		return
	}
	s.sortRange(start, s.len)
	s.dirty = false
}

// sortRange sorts tag entries by ValueOffset for [start,end) using insertion sort.
// The queue is small and fixed-capacity, making this non-alloc strategy effective.
func (s *state) sortRange(start, end uint32) {
	if end-start < 2 {
		return
	}
	for i := start + 1; i < end; i++ {
		v := s.tag[i]
		j := i
		for j > start && v.ValueOffset < s.tag[j-1].ValueOffset {
			s.tag[j] = s.tag[j-1]
			j--
		}
		s.tag[j] = v
	}
}
