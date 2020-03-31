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

const (
	ObjDefault ObjEncoding = iota

	ObjPlainBool
	ObjPlainInt8
	ObjPlainInt16
	ObjPlainInt32
	ObjPlainInt64
	ObjPlainUint8
	ObjPlainUint16
	ObjPlainUint32
	ObjPlainUint64
	ObjPlainFloat32
	ObjPlainFloat64
	ObjPlainString
)

type Object struct {
	Kind     ObjKind
	Encoding ObjEncoding
	Data     interface{}
}

func newObject(kind ObjKind, encoding ObjEncoding, data interface{}) *Object {
	return &Object{
		Kind:     kind,
		Encoding: encoding,
		Data:     data,
	}
}

func PlainObject(any interface{}) *Object {
	// todo
	return newObject(ObjPlain, ObjDefault, any)
}

func ListObject() *Object {
	cls := list.New()
	return newObject(ObjList, ObjDefault, cls)
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
