package exif

import (
	"testing"

	"github.com/evanoberholster/imagemeta/meta/exif/tag"
)

var stateQueueBenchSink uint32

var stateQueueBenchOffsets = [...]uint32{
	92, 11, 204, 17, 88, 399, 77, 133, 51, 290, 12, 410, 67, 309, 5, 222,
	47, 180, 90, 361, 15, 44, 500, 73, 26, 333, 7, 148, 64, 201, 320, 39,
	111, 141, 172, 219, 249, 278, 306, 337, 369, 402, 438, 471, 509, 541, 579, 611,
	18, 29, 42, 56, 71, 87, 104, 122, 143, 165, 188, 212, 237, 263, 291, 318,
}

// oldQueue mirrors the previous insertion-on-add queue behavior.
type oldQueue struct {
	tag [tagQueueMax]tag.Entry
	len uint32
}

func (q *oldQueue) reset() {
	q.len = 0
}

func (q *oldQueue) add(t tag.Entry) bool {
	if q.len >= tagQueueMax {
		return false
	}
	if q.len == 0 {
		q.tag[0] = t
		q.len = 1
		return true
	}
	for i := q.len; i > 0; i-- {
		prev := q.tag[i-1]
		if t.ValueOffset > prev.ValueOffset {
			if i != q.len {
				copy(q.tag[i+1:q.len+1], q.tag[i:q.len])
			}
			q.tag[i] = t
			q.len++
			return true
		}
	}
	copy(q.tag[1:q.len+1], q.tag[0:q.len])
	q.tag[0] = t
	q.len++
	return true
}

func BenchmarkStateQueueBuild(b *testing.B) {
	entries := make([]tag.Entry, len(stateQueueBenchOffsets))
	for i := 0; i < len(stateQueueBenchOffsets); i++ {
		entries[i] = testEntryAtOffset(stateQueueBenchOffsets[i])
	}

	b.Run("OldInsertPerAdd", func(b *testing.B) {
		var q oldQueue
		var sum uint32
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			q.reset()
			for j := 0; j < len(entries); j++ {
				_ = q.add(entries[j])
			}
			for j := uint32(0); j < q.len; j++ {
				sum += q.tag[j].ValueOffset
			}
		}
		stateQueueBenchSink = sum
	})

	b.Run("NewAppendThenSort", func(b *testing.B) {
		var q state
		var sum uint32
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			q.reset()
			for j := 0; j < len(entries); j++ {
				_ = q.addTag(entries[j])
			}
			q.sortAll()
			for j := uint32(0); j < q.len; j++ {
				sum += q.tag[j].ValueOffset
			}
		}
		stateQueueBenchSink = sum
	})
}
