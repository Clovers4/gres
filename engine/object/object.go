package object

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"github.com/clovers4/gres/engine/object/plain"
)

func init() {
	gob.Register(list.List{})
}

type ObjKind byte

const (
	ObjPlain ObjKind = iota
	ObjList
	ObjSet
	ObjZset
	ObjHash
)

var ObjKinds = map[ObjKind]string{
	ObjPlain: "PLAIN", // string, int, float, ...
	ObjList:  "LIST",
	ObjSet:   "SET",
	ObjZset:  "ZSET",
	ObjHash:  "HASH",
}

type ObjEncoding byte

type Object struct {
	Kind ObjKind
	Data interface{}
}

func newObject(kind ObjKind, data interface{}) *Object {
	return &Object{
		Kind: kind,
		Data: data,
	}
}

func PlainObject(any interface{}) *Object {
	return newObject(ObjPlain, any)
}

func ListObject() *Object {
	cls := list.New()
	return newObject(ObjList, cls)
}

func SetObject() *Object {
	m := make(map[string]bool)
	return newObject(ObjSet, m)
}

func HashObject() *Object {
	m := make(map[string]interface{})
	return newObject(ObjHash, m)
}

func (obj *Object) Plain() *plain.Plain {
	return obj.Data.(*plain.Plain)
}

func (obj *Object) List() *list.List {
	if obj.Data != nil {
		return obj.Data.(*list.List)
	}
	return nil
}

func (obj *Object) String() string {
	return fmt.Sprintf("[%v] %v", ObjKinds[obj.Kind], obj.Data)
}
