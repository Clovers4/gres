package proto

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRW(t *testing.T) {
	buf := new(bytes.Buffer)
	var args []interface{}
	var err error
	var val interface{}

	w := NewWriter(buf)
	r := NewReader(buf)

	args = []interface{}{"set", "s", -1, -32.3}
	err = w.ReplyArrays(args)
	assert.Nil(t, err)

	err = w.ReplyStatus("OK")
	assert.Nil(t, err)

	err = w.Flush()
	assert.Nil(t, err)

	val, err = r.ReadReply()
	assert.Nil(t, err)
	assert.Equal(t, "[set s -1 -32.3]", fmt.Sprintf("%v", val))

	val, err = r.ReadReply()
	assert.Nil(t, err)
	assert.Equal(t, "OK", fmt.Sprintf("%v", val))

}
