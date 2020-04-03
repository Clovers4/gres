package set

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet_Marshal(t *testing.T) {
	set := New()
	set.Add("A")
	set.Add("B")
	set.Add(int64(123))
	assert.Equal(t, true, set.Length() > 0)

	// marshal
	buf := new(bytes.Buffer)
	err := set.Marshal(buf)
	assert.Nil(t, err)

	// unmarshal
	newSet := New()
	r := bytes.NewReader(buf.Bytes())
	err = newSet.Unmarshal(r)
	assert.Nil(t, err)

	for val := range set.m {
		_, ok := newSet.m[val]
		fmt.Println(val)
		assert.Equal(t, ok, true)
	}
	assert.Equal(t, set.Length(), newSet.Length())
}
