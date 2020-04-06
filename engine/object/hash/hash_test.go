package hash

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash_Add(t *testing.T) {
	h := New()
	h.Set("A", int64(1))
	h.Set("B", 2.2)
	h.Set("C", "C2")

	fmt.Println(h)
	assert.Equal(t, 3, h.Length())

	h.Delete("A")
	h.Delete("D")
	fmt.Println(h)
	assert.Equal(t, 2, h.Length())

	val, existed := h.Get("B")
	assert.Equal(t, existed, true)
	assert.Equal(t, 2.2, val)
}

func TestHash_Marshal(t *testing.T) {
	h := New()
	h.Set("A", int64(1))
	h.Set("B", 2.2)
	h.Set("C", "C2")

	// marshal
	buf := new(bytes.Buffer)
	err := h.Marshal(buf)
	assert.Nil(t, err)

	// unmarshal
	newH := New()
	r := bytes.NewReader(buf.Bytes())
	err = newH.Unmarshal(r)
	assert.Equal(t, h.Length(), newH.Length())

	fmt.Println(newH)
	assert.Equal(t, h.String(), newH.String())
}
