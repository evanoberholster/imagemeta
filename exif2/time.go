package exif2

import (
	"sync"
	"time"
)

var (
	mapTimeZones   = map[int32]*time.Location{}
	mutexTimeZones = sync.RWMutex{}
)

func getLocation(offset int32, buf []byte) *time.Location {
	mutexTimeZones.RLock()
	if z, ok := mapTimeZones[offset]; ok {
		mutexTimeZones.RUnlock()
		return z
	}
	mutexTimeZones.RUnlock()
	mutexTimeZones.Lock()
	l := time.FixedZone(string(buf), int(offset))
	mapTimeZones[offset] = l
	mutexTimeZones.Unlock()
	return l
}
