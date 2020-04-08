package list

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	ls := New()
	var s1 string
	var s2 string

	// 1
	ls.LPush("A")
	for n := ls.Front(); n != nil; n = n.Next() {
		s1 += fmt.Sprintf("%v ", n.Val())
	}
	for n := ls.End(); n != nil; n = n.Prev() {
		s2 += fmt.Sprintf("%v ", n.Val())
	}
	assert.Equal(t, "{A}", ls.String())
	assert.Equal(t, s1, s2)

	// 2
	ls.LPop()
	s1 = ""
	for n := ls.Front(); n != nil; n = n.Next() {
		s1 += fmt.Sprintf("%v ", n.Val())
	}
	s2 = ""
	for n := ls.End(); n != nil; n = n.Prev() {
		s2 = fmt.Sprintf("%v ", n.Val()) + s2
	}
	assert.Equal(t, s1, s2)
	assert.Equal(t, 0, ls.Length())

	// 3
	ls.LPop()
	assert.Equal(t, 0, ls.Length())

	// 4
	ls.RPush("B")
	ls.RPush("C")
	ls.LPush("A")
	s1 = ""
	for n := ls.Front(); n != nil; n = n.Next() {
		s1 += fmt.Sprintf("%v ", n.Val())
	}
	s2 = ""
	for n := ls.End(); n != nil; n = n.Prev() {
		s2 = fmt.Sprintf("%v ", n.Val()) + s2
	}
	assert.Equal(t, s1, s2)
	assert.Equal(t, "{A, B, C}", ls.String())
	assert.Equal(t, "A", ls.Index(0).Val())
	assert.Equal(t, "B", ls.Index(1).Val())
	assert.Equal(t, "C", ls.Index(2).Val())
	assert.Equal(t, "C", ls.Index(-1).Val())
	assert.Nil(t, ls.Index(3))
	assert.Nil(t, ls.Index(-4))

	// 5
	ls.LPop()
	ls.RPop()
	s1 = ""
	for n := ls.Front(); n != nil; n = n.Next() {
		s1 += fmt.Sprintf("%v ", n.Val())
	}
	s2 = ""
	for n := ls.End(); n != nil; n = n.Prev() {
		s2 = fmt.Sprintf("%v ", n.Val()) + s2
	}
	assert.Equal(t, s1, s2)
	assert.Equal(t, "{B}", ls.String())
}

func TestList_Marshal(t *testing.T) {
	ls := New()
	ls.LPush("A")
	ls.LPush("B")
	ls.LPush(int64(13))
	ls.LPush(3.5)

	// marshal
	buf := new(bytes.Buffer)
	err := ls.Marshal(buf)
	assert.Nil(t, err)

	// unmarshal
	newLs := New()
	r := bytes.NewReader(buf.Bytes())
	err = newLs.Unmarshal(r)
	assert.Equal(t, ls.Length(), newLs.Length())
	assert.Equal(t, ls.String(), newLs.String())
	fmt.Println(newLs)

	for n, nn := ls.Front(), newLs.Front(); n != nil && nn != nil; n, nn = n.Next(), nn.Next() {
		fmt.Println(n.Val(), nn.Val())
		assert.Equal(t, n.Val(), nn.Val())
	}

}
