package object

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObject_Marshal_Plain(t *testing.T) {
	obj := PlainObject(int64(12))
	fmt.Println(obj)

	// marshal
	buf := new(bytes.Buffer)
	err := obj.Marshal(buf)
	assert.Nil(t, err)

	// unmarshal
	newObj := &Object{}
	r := bytes.NewReader(buf.Bytes())
	err = newObj.Unmarshal(r)
	assert.Equal(t, obj.String(), newObj.String())
	fmt.Println(newObj.String())
}

func TestObject_Marshal_List(t *testing.T) {
	obj := ListObject()
	ls, ok := obj.List()
	assert.Equal(t, true, ok)
	ls.RPush("A")
	ls.RPush("B")
	fmt.Println(obj)

	// marshal
	buf := new(bytes.Buffer)
	err := obj.Marshal(buf)
	assert.Nil(t, err)

	// unmarshal
	newObj := &Object{}
	r := bytes.NewReader(buf.Bytes())
	err = newObj.Unmarshal(r)
	assert.Equal(t, obj.String(), newObj.String())
	fmt.Println(newObj.String())
}

func TestObject_Marshal_Set(t *testing.T) {
	obj := SetObject()
	set, ok := obj.Set()
	assert.Equal(t, true, ok)
	set.Add("A")
	set.Add("B")
	fmt.Println(obj)

	// marshal
	buf := new(bytes.Buffer)
	err := obj.Marshal(buf)
	assert.Nil(t, err)

	// unmarshal
	newObj := &Object{}
	r := bytes.NewReader(buf.Bytes())
	err = newObj.Unmarshal(r)
	assert.Equal(t, obj.String(), newObj.String())
	fmt.Println(newObj.String())
}

func TestObject_Marshal_ZSet(t *testing.T) {
	obj := ZSetObject()
	set, ok := obj.ZSet()
	assert.Equal(t, true, ok)
	set.Add(23.0, "A")
	set.Add(12, "B")
	fmt.Println(obj)

	// marshal
	buf := new(bytes.Buffer)
	err := obj.Marshal(buf)
	assert.Nil(t, err)

	// unmarshal
	newObj := &Object{}
	r := bytes.NewReader(buf.Bytes())
	err = newObj.Unmarshal(r)
	assert.Equal(t, obj.String(), newObj.String())
	fmt.Println(newObj.String())
}

func TestObject_Marshal_Hash(t *testing.T) {
	obj := HashObject()
	set, ok := obj.Hash()
	assert.Equal(t, true, ok)
	set.Add("A", int64(23))
	set.Add("B", "B23")
	fmt.Println(obj)

	_, ok = obj.List()
	assert.Equal(t, false, ok)

	// marshal
	buf := new(bytes.Buffer)
	err := obj.Marshal(buf)
	assert.Nil(t, err)

	// unmarshal
	newObj := &Object{}
	r := bytes.NewReader(buf.Bytes())
	err = newObj.Unmarshal(r)
	assert.Equal(t, obj.String(), newObj.String())
	fmt.Println(newObj.String())
}
