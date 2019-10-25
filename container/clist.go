package container

import (
	"container/list"
	"sync"
)

type CList struct {
	l  *list.List
	mu sync.RWMutex
}

func NewCList() *CList {
	return &CList{
		l: list.New(),
	}
}

func (cls *CList) Len() int {
	cls.mu.RLock()
	len := cls.l.Len()
	cls.mu.RUnlock()
	return len
}

func (cls *CList) LPush(value ...string) int {
	cls.mu.Lock()
	for _, v := range value {
		cls.l.PushFront(v)
	}
	len := cls.l.Len()
	cls.mu.Unlock()
	return len
}

func (cls *CList) LPop() string {
	cls.mu.Lock()
	old := cls.l.Remove(cls.l.Front())
	cls.mu.Unlock()
	return old.(string)
}

func (cls *CList) RPush(value ...interface{}) int {
	cls.mu.Lock()
	for _, v := range value {
		cls.l.PushBack(v)
	}
	len := cls.l.Len()
	cls.mu.Unlock()
	return len
}

func (cls *CList) RPop() string {
	cls.mu.Lock()
	old := cls.l.Remove(cls.l.Back())
	cls.mu.Unlock()
	return old.(string)
}

func (cls *CList) Range(start, stop int) []string {
	if stop < 0 {
		stop = cls.l.Len() + stop
	}

	ls := make([]string, 0, stop-start+1)
	cls.mu.RLock()
	for now, e := 0, cls.l.Front(); now < cls.Len() && now <= stop && e != nil; now, e = now+1, e.Next() {
		if now < start {
			continue
		}
		ls = append(ls, e.Value.(string))
	}
	cls.mu.RUnlock()
	return ls
}
