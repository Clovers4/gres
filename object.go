package gres

import (
	"fmt"
	"github.com/clovers4/gres/container"
)

const (
	OBJ_STRING = 0
	OBJ_LIST   = 1
	OBJ_SET    = 2
	OBJ_ZSET   = 3
	OBJ_HASH   = 4
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

func (obj *object) getString() (string, error) {
	if obj.kind != OBJ_STRING {
		return "", ErrWrongTypeOps
	}
	return obj.ptr.(string), nil
}

func (obj *object) getList() (*container.CList, error) {
	if obj.kind != OBJ_LIST {
		return nil, ErrWrongTypeOps
	}
	return obj.ptr.(*container.CList), nil
}

func (obj *object) String() string {
	return fmt.Sprintf("%v", obj.ptr)
}
