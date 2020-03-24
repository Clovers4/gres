package gres

import "github.com/clovers4/gres/container"

type DB struct {
	all *container.CMap
	//todo
	save *container.CMap
}

func NewDB() *DB {
	return &DB{
		all: container.NewCMap(),
	}
}

func (db *DB) Set(key string, obj *object) *object {
	if v := db.all.Set(key, obj); v != nil {
		return v.(*object)
	}
	return nil
}

func (db *DB) Get(key string) *object {
	v, ok := db.all.Get(key)
	if !ok {
		return nil
	}
	return v.(*object)
}

func (db *DB) CheckKind(key string, kind int) bool {
	obj := db.Get(key)
	if obj == nil {
		return true
	}
	if obj.kind == kind {
		return true
	}
	return false
}
