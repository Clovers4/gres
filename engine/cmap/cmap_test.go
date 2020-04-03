package cmap

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Animal struct {
	name string
}

func TestNewCMap(t *testing.T) {
	cm := New()
	assert.NotNil(t, cm)
}

func TestCMap_Set_Get(t *testing.T) {
	cm := New()
	elephant := Animal{"elephant"}
	cm.Set("elephant", elephant)

	e, oke := cm.Get(elephant.name)
	assert.Equal(t, true, oke)
	assert.Equal(t, e, elephant)

	v, okv := cm.Get("not found")
	assert.Equal(t, false, okv)
	assert.Equal(t, nil, v)
}

func TestCMap_Exist(t *testing.T) {
	cm := New()
	elephant := Animal{"elephant"}
	cm.Set("elephant", elephant)

	assert.Equal(t, true, cm.Exist(elephant.name))
	assert.Equal(t, false, cm.Exist("not found"))
}

func TestCMap_Remove(t *testing.T) {
	cm := New()
	elephant := Animal{"elephant"}
	cm.Set("elephant", elephant)
	assert.Equal(t, true, cm.Exist(elephant.name))

	cm.Remove(elephant.name)
	assert.Equal(t, false, cm.Exist(elephant.name))
}

func TestCMap_Count_ConcurrentSuccess(t *testing.T) {
	loopN := 100
	cm := New()
	start := sync.WaitGroup{}
	start.Add(1)
	end := sync.WaitGroup{}
	end.Add(loopN)
	for i := 0; i < loopN; i++ {
		go func() {
			start.Wait()
			cm.Set(strconv.Itoa(i), i)
			end.Done()
		}()
	}
	start.Done()
	end.Wait()
	assert.NotEqual(t, loopN, cm.Count())
}

// TestConcurrentFail will result in -- fatal error: concurrent map writes
// It should be ignored, because it cannot be caught by recover()
//
//func TestConcurrentFail(t *testing.T) {
//	loopN := 100
//	items := make(map[string]int)
//	start := sync.WaitGroup{}
//	start.Add(1)
//	end := sync.WaitGroup{}
//	end.Add(loopN)
//	for i := 0; i < loopN; i++ {
//		go func() {
//			start.Wait()
//			items[strconv.Itoa(i)] = i
//			end.Done()
//		}()
//	}
//	start.Done()
//	end.Wait()
//	fmt.Println(len(items),items)
//	assert.NotEqual(t, loopN, len(items))
//}
//
//func TestCMap_Marshal(t *testing.T) {
//	cm := New()
//	cm.Set("plain-A", object.PlainObject("A-and"))
//	cm.Set("plain-B", object.PlainObject("B-bb"))
//
//	lsObj := object.ListObject()
//	ls, _ := lsObj.List()
//	ls.RPush("C1")
//	ls.RPush("C2")
//	cm.Set("list-C", lsObj)
//
//	setObj := object.SetObject()
//	set, _ := setObj.Set()
//	set.Add("A")
//	set.Add("B")
//	cm.Set("set-D", setObj)
//
//	zsetObj := object.ZSetObject()
//	zs, _ := zsetObj.ZSet()
//	zs.Add(23, "ZS_A")
//	zs.Add(23, "ZS_B")
//	zs.Add(3, "ZS_C")
//	cm.Set("zset-E", zsetObj)
//
//	hashObj := object.HashObject()
//	ha, _ := hashObj.Hash()
//	ha.Add("a", "b")
//	ha.Add("b", 2)
//	cm.Set("hash-F", hashObj)
//
//	fmt.Println(cm)
//
//}
