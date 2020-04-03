package cmap

import (
	"fmt"
	fnv2 "hash/fnv"
	"sort"
	"sync"
)

const defaultSegmentCount = 32

// todo: 动态调整
type cmapSegment struct {
	items        map[string]interface{} // the value must be object
	sync.RWMutex                        // Read Write mutex, guards access to internal map.
}

type CMap struct {
	segments []*cmapSegment
}

func New() *CMap {
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

func hash(key string) uint32 {
	h := fnv2.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

// Only for test
func (cm *CMap) String() string {
	var s string
	s += "{"

	var keys []string
	for _, seg := range cm.segments {
		seg.RLock()
		defer seg.RUnlock()
		for k := range seg.items {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, key := range keys {
		val, _ := cm.Get(key)
		s += fmt.Sprintf("%v : %v, ", key, val)
	}

	s = s[:len(s)-2]
	s += "}"
	return s
}

//
//// 实际上整个 Marshal 过程 cmap 都需要被锁定 or 保证只有一个
//func (cm *CMap) Marshal(w io.Writer) error {
//	// write total.
//	total := cm.Count()
//	if err := util.Write(w, int64(total)); err != nil {
//		return err
//	}
//
//	// loop write score and val
//	for _, seg := range cm.segments {
//		seg.RLock()
//		defer seg.RUnlock()
//		for key, obj := range seg.items {
//			if err := util.Write(w, key); err != nil {
//				return err
//			}
//
//			v := obj.(serialize.Serializable)
//			if err := v.Marshal(w); err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}
//
//func (cm *CMap) Unmarshal(r io.Reader) error {
//	var total int64
//	if err := util.Read(r, &total); err != nil {
//		return err
//	}
//
//	for i := 0; i < int(total); i++ {
//		var val string
//		if err := util.Read(r, &val); err != nil {
//			return err
//		}
//
//		obj := object.NilObject()
//		if err := obj.Unmarshal(r); err != nil {
//			return err
//		}
//
//		cm.Set(val, obj)
//	}
//	return nil
//}
