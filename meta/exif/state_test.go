package exif

import (
	"testing"

	"github.com/evanoberholster/imagemeta/meta/exif/ifd"
	"github.com/evanoberholster/imagemeta/meta/exif/tag"
	"github.com/evanoberholster/imagemeta/meta/utils"
)

func testEntryAtOffset(offset uint32) tag.Entry {
	return tag.NewEntry(tag.TagMake, tag.TypeLong, 1, offset, ifd.IFD0, 0, utils.LittleEndian)
}

func TestStateAddTagOrdersByOffsetAfterSort(t *testing.T) {
	t.Parallel()

	var s state
	if !s.addTag(testEntryAtOffset(30)) || !s.addTag(testEntryAtOffset(10)) || !s.addTag(testEntryAtOffset(20)) {
		t.Fatal("addTag() returned false unexpectedly")
	}

	if s.len != 3 {
		t.Fatalf("len = %d, want 3", s.len)
	}
	if !s.dirty {
		t.Fatal("queue should be marked dirty before sorting")
	}
	s.sortAll()
	if s.tag[0].ValueOffset != 10 || s.tag[1].ValueOffset != 20 || s.tag[2].ValueOffset != 30 {
		t.Fatalf("tag order = [%d,%d,%d], want [10,20,30]", s.tag[0].ValueOffset, s.tag[1].ValueOffset, s.tag[2].ValueOffset)
	}
}

func TestStateTagNavigationAndCompaction(t *testing.T) {
	t.Parallel()

	var s state
	_ = s.addTag(testEntryAtOffset(10))
	_ = s.addTag(testEntryAtOffset(20))
	_ = s.addTag(testEntryAtOffset(30))

	if !s.validTag() {
		t.Fatal("validTag() should be true with queued tags")
	}
	if got := s.currentTag().ValueOffset; got != 10 {
		t.Fatalf("currentTag() = %d, want 10", got)
	}
	if got := s.advanceTag().ValueOffset; got != 20 {
		t.Fatalf("advanceTag() first = %d, want 20", got)
	}
	if got := s.advanceTag().ValueOffset; got != 30 {
		t.Fatalf("advanceTag() second = %d, want 30", got)
	}
	if got := s.advanceTag(); got != (tag.Entry{}) {
		t.Fatalf("advanceTag() end = %+v, want zero entry", got)
	}
	if s.validTag() {
		t.Fatal("validTag() should be false after exhausting queue")
	}

	s.pos = 1
	s.resetPosition()
	if s.pos != 0 {
		t.Fatalf("resetPosition() pos = %d, want 0", s.pos)
	}
	if s.len != 2 {
		t.Fatalf("resetPosition() len = %d, want 2", s.len)
	}
	if s.tag[0].ValueOffset != 20 || s.tag[1].ValueOffset != 30 {
		t.Fatalf("resetPosition() queue head = [%d,%d], want [20,30]", s.tag[0].ValueOffset, s.tag[1].ValueOffset)
	}
}

func TestStateResetAndQueueFull(t *testing.T) {
	t.Parallel()

	var s state
	s.len = 7
	s.pos = 3
	s.reset()
	if s.len != 0 || s.pos != 0 {
		t.Fatalf("reset() = (len=%d,pos=%d), want (0,0)", s.len, s.pos)
	}

	s.len = tagQueueMax
	if s.addTag(testEntryAtOffset(1)) {
		t.Fatal("addTag() should fail when queue is full")
	}
}
