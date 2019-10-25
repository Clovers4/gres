package container

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
	cm := NewCMap()
	assert.NotNil(t, cm)
}

func TestCMap_Set_Get(t *testing.T) {
	cm := NewCMap()
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
	cm := NewCMap()
	elephant := Animal{"elephant"}
	cm.Set("elephant", elephant)

	assert.Equal(t, true, cm.Exist(elephant.name))
	assert.Equal(t, false, cm.Exist("not found"))
}

func TestCMap_Remove(t *testing.T) {
	cm := NewCMap()
	elephant := Animal{"elephant"}
	cm.Set("elephant", elephant)
	assert.Equal(t, true, cm.Exist(elephant.name))

	cm.Remove(elephant.name)
	assert.Equal(t, false, cm.Exist(elephant.name))
}

func TestCMap_Count_ConcurrentSuccess(t *testing.T) {
	loopN := 100
	cm := NewCMap()
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
//	m := make(map[string]int)
//	start := sync.WaitGroup{}
//	start.Add(1)
//	end := sync.WaitGroup{}
//	end.Add(loopN)
//	for i := 0; i < loopN; i++ {
//		go func() {
//			start.Wait()
//			m[strconv.Itoa(i)] = i
//			end.Done()
//		}()
//	}
//	start.Done()
//	end.Wait()
//	fmt.Println(len(m),m)
//	assert.NotEqual(t, loopN, len(m))
//}
