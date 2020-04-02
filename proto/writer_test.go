package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/clovers4/gres/util"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func dos(b []byte) {
	b = b[4:]
	fmt.Println("dos", b)
}

func TestByte(t *testing.T) {
	b := []byte{1, 2, 3, 4, 5, 6, 7, 78}
	fmt.Println(b)
	dos(b)
	fmt.Println(b)

}

func TestWriteByBinary(t *testing.T) {
	var err error
	buf := &bytes.Buffer{}

	s1 := "asdasd"
	b1 := util.StringToBytes(s1)
	num := int64(2232331)

	fmt.Println(b1)

	err = binary.Write(buf, binary.BigEndian, b1)
	assert.Nil(t, err)
	fmt.Println(buf.Bytes())

	err = binary.Write(buf, binary.BigEndian, num)
	assert.Nil(t, err)
	fmt.Println(buf.Bytes())

	var s2 string
	var b2 []byte
	var num2 int64
	fmt.Println(buf.Bytes()[:len(b1)])
	b := bytes.NewReader(buf.Bytes()[:len(b1)])
	err = binary.Read(b, binary.LittleEndian, &b2)
	fmt.Println("b2", b2)
	assert.Nil(t, err)
	fmt.Println(buf.Bytes())

	err = binary.Read(b, binary.LittleEndian, &num2)
	assert.Nil(t, err)
	fmt.Println(buf.Bytes())

	s2 = util.BytesToString(b2)
	fmt.Println(s1, s2)
	fmt.Println(num, num2)

}

func TestWriteByString(t *testing.T) {
	buf := []byte{}
	n := int64(1244444232343)
	buf = strconv.AppendInt(buf[:0], n, 10)
	fmt.Println(buf)
	fmt.Println(len(buf))

}
func TestWrite(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewWriter(buf)
	w.string("ABC")
	w.Flush()

	fmt.Println(buf.Bytes())
}

var startSeed = time.Now().UnixNano()

func randomSlice() []byte {
	b := make([]byte, 0, rand.Intn(100))
	for i, v := range b {
		b[i] = v
	}
	return b
}

func BenchmarkAppend(b *testing.B) {
	rand.Seed(startSeed)
	b.ResetTimer()
	var all []byte

	for i := 0; i < b.N; i++ {
		all = append(all, randomSlice()...)
	}
}

func BenchmarkBufferWrite(b *testing.B) {
	rand.Seed(startSeed)
	b.ResetTimer()
	var buff bytes.Buffer
	for i := 0; i < b.N; i++ {
		buff.Write(randomSlice())
	}
	all := buff.Bytes()
	_ = all
}
