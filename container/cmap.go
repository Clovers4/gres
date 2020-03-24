package container

import (
	fnv2 "hash/fnv"
	"sync"
)

var defaultSegmentCount = 32

// todo:动态调整
type cmapSegment struct {
	m  map[string]interface{}
	mu sync.RWMutex // Read Write mutex, guards access to internal map.
}

type CMap struct {
	segs []*cmapSegment
}

func NewCMap() *CMap {
	m := &CMap{
		segs: make([]*cmapSegment, defaultSegmentCount),
	}
	for i := 0; i < defaultSegmentCount; i++ {
		m.segs[i] = &cmapSegment{m: make(map[string]interface{})}
	}
	return m
}

// getSegment returns segment related the given key
func (cm *CMap) getSegment(key string) *cmapSegment {
	return cm.segs[hash(key)%uint32(defaultSegmentCount)]
}

// Set the given value under the specified key.
func (cm *CMap) Set(key string, value interface{}) interface{}{
	seg := cm.getSegment(key)
	seg.mu.Lock()
	old:=seg.m[key]
	seg.m[key] = value
	seg.mu.Unlock()
	return old
}

// Get retrieves an element from map under given key.
func (cm *CMap) Get(key string) (interface{}, bool) {
	seg := cm.getSegment(key)
	seg.mu.RLock()
	val, ok := seg.m[key]
	seg.mu.RUnlock()
	return val, ok
}

func (cm *CMap) Exist(key string) bool {
	seg := cm.getSegment(key)
	seg.mu.RLock()
	_, ok := seg.m[key]
	seg.mu.RUnlock()
	return ok
}

// Remove removes an element from the map.
func (cm *CMap) Remove(key string) (v interface{}, ok bool) {
	seg := cm.getSegment(key)
	seg.mu.Lock()
	v, ok = seg.m[key]
	delete(seg.m, key)
	seg.mu.Unlock()
	return v, ok
}

// Count returns amount of elements in CMap.
// But the count is not very accurate.
func (cm *CMap) Count() int {
	c := 0
	for _, seg := range cm.segs {
		seg.mu.RLock()
		c += len(seg.m)
		seg.mu.RUnlock()
	}
	return c
}

func hash(key string) uint32 {
	h := fnv2.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}
