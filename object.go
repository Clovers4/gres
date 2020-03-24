package gres

import (
	"fmt"
	"github.com/clovers4/gres/container"
)

const (
	OBJ_STRING = iota
	OBJ_LIST
	OBJ_SET
	OBJ_ZSET
	OBJ_HASH
)

type object struct {
	kind int
	ptr  interface{}
}

func createObject(kind int, ptr interface{}) *object {
	return &object{
		kind: kind,
		ptr:  ptr,
	}
}

func createStringObject(s string) *object {
	return createObject(OBJ_STRING, s)
}

func createListObject() *object {
	cls := container.NewCList()
	return createObject(OBJ_LIST, cls)
}

func (obj *object) checkKind(kind int) bool {
	return obj.kind == kind
}

func (obj *object) getString() string {
	if obj.kind != OBJ_STRING {
		panic(ErrWrongTypeOps)
	}
	if obj.ptr!=nil{
		return obj.ptr.(string)
	}
	return ""
}

func (obj *object) getList() *container.CList {
	if obj.kind != OBJ_LIST {
		panic(ErrWrongTypeOps)
	}
	return obj.ptr.(*container.CList)
}

func (obj *object) String() string {
	return fmt.Sprintf("%v", obj.ptr)
}
