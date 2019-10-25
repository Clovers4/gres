package gres

import "github.com/clovers4/gres/container"

type DB struct {
	id  int // Database ID
	all *container.CMap
	srv *Server
}

func NewDB(id int, srv *Server) *DB {
	return &DB{
		id:  id,
		srv: srv,
		all: container.NewCMap(),
	}
}

func (db *DB) Set(key string, obj *object) *object {
	if v := db.all.Set(key, obj); v != nil {
		return v.(*object)
	}
	return nil
}

func (db *DB) Get(key string) (obj *object, ok bool) {
	v, ok := db.all.Get(key)
	if !ok {
		return nil, false
	}
	return v.(*object), ok
}
