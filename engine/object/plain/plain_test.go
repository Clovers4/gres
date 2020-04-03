package plain

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlain_Marshal_String(t *testing.T) {
	p := New("marshal is success")

	// marshal
	buf := new(bytes.Buffer)
	err := p.Marshal(buf)
	assert.Nil(t, err)

	// unmarshal
	newP := New(nil)
	r := bytes.NewReader(buf.Bytes())
	err = newP.Unmarshal(r)
	assert.Equal(t, p.Val(), newP.Val())
}

func TestPlain_Marshal_Int64(t *testing.T) {
	p := New(int64(23))

	// marshal
	buf := new(bytes.Buffer)
	err := p.Marshal(buf)
	assert.Nil(t, err)

	// unmarshal
	newP := New(nil)
	r := bytes.NewReader(buf.Bytes())
	err = newP.Unmarshal(r)
	assert.Equal(t, p.Val(), newP.Val())
}

func TestPlain_Marshal_Float64(t *testing.T) {
	p := New(float64(23.2333))

	// marshal
	buf := new(bytes.Buffer)
	err := p.Marshal(buf)
	assert.Nil(t, err)

	// unmarshal
	newP := New(nil)
	r := bytes.NewReader(buf.Bytes())
	err = newP.Unmarshal(r)
	assert.Equal(t, p.Val(), newP.Val())
}
