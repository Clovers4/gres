package object

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestString2File(t *testing.T) {
	//newFile, err := os.OpenFile("string_f", os.O_CREATE, 0666)
	//s := "GRES"
	//origin := PlainObject(s)
	//data, err := MarshalObj(origin)
}

func TestString(t *testing.T) {
	s := "GRES"
	origin := PlainObject(s)
	data, err := MarshalObj(origin)
	assert.Nil(t, err)

	derived, err := UnmarshalObj(data)
	assert.Equal(t, s, derived.Plain())
	fmt.Println(derived.Plain())
}

func TestString2(t *testing.T) {
	s := "GRES"
	origin := PlainObject(s)
	data, err := MarshalObjString(origin)
	assert.Nil(t, err)

	derived, err := UnmarshalObjString(data)
	assert.Equal(t, s, derived.Plain())
	fmt.Println(derived.Plain())
}

func TestList(t *testing.T) {
	origin := ListObject()
	list := origin.List()
	list.PushBack("1")
	list.PushBack("2")
	fmt.Println(origin)

	data, err := MarshalObj(origin)
	assert.Nil(t, err)

	derived, err := UnmarshalObj(data)
	fmt.Println(derived)
}

func TestList2(t *testing.T) {
	buf := strconv.AppendInt([]byte{}, 123, 10)
	fmt.Println(buf)
}
