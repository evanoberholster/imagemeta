package exif2

import (
	"sync"
	"time"
)

var (
	// cacheTimeZone caches time.Location to avoid allocs and increase performance.
	// time.Location should only need to be calculated once.
	cacheTimeZone  = map[int32]*time.Location{}
	mutexTimeZones = sync.RWMutex{}
)

// getLocation faciliates an offset and a time string to result
// with a *time.Location creating it when not cound in the cache.
// RWMutex for concurrancy.
func getLocation(offset int32, buf []byte) *time.Location {
	mutexTimeZones.RLock()
	if z, ok := cacheTimeZone[offset]; ok {
		mutexTimeZones.RUnlock()
		return z
	}
	mutexTimeZones.RUnlock()
	mutexTimeZones.Lock()
	l := time.FixedZone(string(buf), int(offset))
	cacheTimeZone[offset] = l
	mutexTimeZones.Unlock()
	return l
}
