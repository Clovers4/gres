package container

import (
	fnv2 "hash/fnv"
	"sync"
)

const defaultSegmentCount = 32

// todo: 动态调整
type cmapSegment struct {
	items        map[string]interface{}
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

type CMap struct {
	segments []*cmapSegment
}

func NewCMap() *CMap {
	cm := &CMap{
		segments: make([]*cmapSegment, defaultSegmentCount),
	}
	for i := 0; i < defaultSegmentCount; i++ {
		cm.segments[i] = &cmapSegment{items: make(map[string]interface{})}
	}
	return cm
}

// getSegment returns segment related the given key
func (cm *CMap) getSegment(key string) *cmapSegment {
	return cm.segments[hash(key)%uint32(defaultSegmentCount)]
}

// set the given value under the specified key.
func (cm *CMap) Set(key string, value interface{}) (val interface{}, existed bool) {
	seg := cm.getSegment(key)
	seg.Lock()
	old, ok := seg.items[key]
	seg.items[key] = value
	seg.Unlock()
	return old, ok
}

// get retrieves an element from map under given key.
func (cm *CMap) Get(key string) (val interface{}, existed bool) {
	seg := cm.getSegment(key)
	seg.RLock()
	val, ok := seg.items[key]
	seg.RUnlock()
	return val, ok
}

func (cm *CMap) Exist(key string) bool {
	seg := cm.getSegment(key)
	seg.RLock()
	_, ok := seg.items[key]
	seg.RUnlock()
	return ok
}

// Remove removes an element from the map.
func (cm *CMap) Remove(key string) (v interface{}, ok bool) {
	seg := cm.getSegment(key)
	seg.Lock()
	v, ok = seg.items[key]
	delete(seg.items, key)
	seg.Unlock()
	return v, ok
}

// Count returns amount of elements in CMap.
// But the count is not very accurate.
func (cm *CMap) Count() int {
	count := 0
	for _, seg := range cm.segments {
		seg.RLock()
		count += len(seg.items)
		seg.RUnlock()
	}
	return count
}

type KVPair struct {
	Key   string
	Value interface{}
}

// Snapshot cannot get real consistent data map expect in locked.
func (cm *CMap) Snapshot() (int, <-chan KVPair) {
	count := cm.Count()
	ch := make(chan KVPair)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(cm.segments))

		for _, seg := range cm.segments {
			go func(seg *cmapSegment) {
				seg.RLock()
				for key, val := range seg.items {
					ch <- KVPair{key, val}
				}
				seg.RUnlock()
				wg.Done()
			}(seg)
		}
		wg.Wait()
		close(ch) //todo:?
	}()
	return count, ch
}

func hash(key string) uint32 {
	h := fnv2.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}
