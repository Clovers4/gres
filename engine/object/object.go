package object

import (
	"fmt"
	"io"

	"github.com/clovers4/gres/engine/object/hash"
	"github.com/clovers4/gres/engine/object/list"
	"github.com/clovers4/gres/engine/object/plain"
	"github.com/clovers4/gres/engine/object/set"
	"github.com/clovers4/gres/engine/object/zset"
	"github.com/clovers4/gres/serialize"
	"github.com/clovers4/gres/util"
)

var Expunged = &Object{}

type ObjKind uint8

func (kind ObjKind) String() string {
	return ObjKinds[kind]
}

const (
	ObjPlain ObjKind = iota
	ObjList
	ObjSet
	ObjZset
	ObjHash
)

var ObjKinds = map[ObjKind]string{
	ObjPlain: "PLAIN", // string, int, float, ... todo: 兼容redis协议
	ObjList:  "LIST",
	ObjSet:   "SET",
	ObjZset:  "ZSET",
	ObjHash:  "HASH",
}

type Object struct {
	kind ObjKind
	data interface{}
}

func newObject(kind ObjKind, data interface{}) *Object {
	return &Object{
		kind: kind,
		data: data,
	}
}

func PlainObject(val interface{}) *Object {
	return newObject(ObjPlain, plain.New(val))
}

func ListObject() *Object {
	return newObject(ObjList, list.New())
}

func SetObject() *Object {
	return newObject(ObjSet, set.New())
}

func ZSetObject() *Object {
	return newObject(ObjZset, zset.New())
}

func HashObject() *Object {
	return newObject(ObjHash, hash.New())
}

func (obj *Object) Kind() ObjKind {
	return obj.kind
}

func (obj *Object) Plain() (*plain.Plain, bool) {
	p, ok := obj.data.(*plain.Plain)
	return p, ok
}

func (obj *Object) List() (*list.List, bool) {
	ls, ok := obj.data.(*list.List)
	return ls, ok
}

func (obj *Object) Set() (*set.Set, bool) {
	s, ok := obj.data.(*set.Set)
	return s, ok
}

func (obj *Object) ZSet() (*zset.ZSet, bool) {
	zs, ok := obj.data.(*zset.ZSet)
	return zs, ok
}

func (obj *Object) Hash() (*hash.Hash, bool) {
	h, ok := obj.data.(*hash.Hash)
	return h, ok
}

func (obj *Object) String() string {
	return fmt.Sprintf("[%v] %v", ObjKinds[obj.kind], obj.data)
}

func (obj *Object) Marshal(w io.Writer) error {
	kind := uint8(obj.kind)
	if err := util.Write(w, kind); err != nil {
		return err
	}

	data := obj.data.(serialize.Serializable)
	if err := data.Marshal(w); err != nil {
		return err
	}
	return nil
}

func (obj *Object) Unmarshal(r io.Reader) error {
	var kind uint8
	if err := util.Read(r, &kind); err != nil {
		return err
	}
	obj.kind = ObjKind(kind)

	switch ObjKind(kind) {
	case ObjPlain:
		data := plain.New(nil)
		if err := data.Unmarshal(r); err != nil {
			return err
		}
		obj.data = data
	case ObjList:
		data := list.New()
		if err := data.Unmarshal(r); err != nil {
			return err
		}
		obj.data = data
	case ObjSet:
		data := set.New()
		if err := data.Unmarshal(r); err != nil {
			return err
		}
		obj.data = data
	case ObjZset:
		data := zset.New()
		if err := data.Unmarshal(r); err != nil {
			return err
		}
		obj.data = data
	case ObjHash:
		data := hash.New()
		if err := data.Unmarshal(r); err != nil {
			return err
		}
		obj.data = data
	default:
		return fmt.Errorf("unsupported object type [%v]", kind)
	}
	return nil
}
