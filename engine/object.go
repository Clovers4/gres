package engine

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"encoding/gob"
	"fmt"
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

func (obj *Object) Plain() string {
	if obj.Data != nil {
		return fmt.Sprintf("%v", obj.Data)
	}
	return ""
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

func MarshalObj(obj *Object) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MarshalObjString(obj *Object) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, obj.Data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func UnmarshalObj(b []byte) (*Object, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)

	obj := &Object{}
	err := decoder.Decode(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func UnmarshalObjString(b []byte) (*Object, error) {
	r := bytes.NewReader(b)
	obj := &Object{}
	if err := binary.Read(r, binary.LittleEndian, &obj.Data); err != nil {
		return nil, err

	}
	return obj, nil
}
