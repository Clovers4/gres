package list

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	ls := New()

	// 1
	ls.LPush("A")
	for n := ls.Front(); n != nil; n = n.Next() {
		fmt.Print(n.Val(), " ")
	}
	fmt.Print(": ")
	for n := ls.End(); n != nil; n = n.Prev() {
		fmt.Print(n.Val(), " ")
	}
	fmt.Println()

	// 2
	ls.LPop()
	for n := ls.Front(); n != nil; n = n.Next() {
		fmt.Print(n.Val(), " ")
	}
	fmt.Print(": ")
	for n := ls.End(); n != nil; n = n.Prev() {
		fmt.Print(n.Val(), " ")
	}
	fmt.Println()

	// 3
	ls.LPop()
	assert.Equal(t, 0, ls.Length())

	// 4
	ls.RPush("B")
	ls.RPush("C")
	ls.LPush("A")
	for n := ls.Front(); n != nil; n = n.Next() {
		fmt.Print(n.Val(), " ")
	}
	fmt.Print(": ")
	for n := ls.End(); n != nil; n = n.Prev() {
		fmt.Print(n.Val(), " ")
	}
	fmt.Println()

	// 5
	ls.LPop()
	ls.RPop()
	for n := ls.Front(); n != nil; n = n.Next() {
		fmt.Print(n.Val(), " ")
	}
	fmt.Print(": ")
	for n := ls.End(); n != nil; n = n.Prev() {
		fmt.Print(n.Val(), " ")
	}
	fmt.Println()
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
