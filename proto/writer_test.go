package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	buf := &bytes.Buffer{}

	num := "asdasd"
	err := binary.Write(buf, binary.LittleEndian, num)
	assert.Nil(t, err)
	fmt.Println(buf.Bytes())

	var num2 string
	b := bytes.NewReader(buf.Bytes())
	binary.Read(b, binary.LittleEndian, &num2)
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
