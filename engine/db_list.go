package engine

import (
	"fmt"
	"github.com/clovers4/gres/engine/object"
	"github.com/clovers4/gres/engine/object/list"
)

// todo:list hash   num shrink

// ========
//   List
// ========
func (db *DB) LPush(key string, val ...interface{}) (int, error) {
	obj := db.get(key)
	if obj == nil {
		obj = object.ListObject()
		db.set(key, obj)
	}

	ls, ok := obj.List()
	if !ok {
		return 0, ErrWrongTypeOps
	}

	for _, v := range val {
		ls.LPush(v)
	}
	return ls.Length(), nil
}

func (db *DB) RPush(key string, val ...interface{}) (int, error) {
	obj := db.get(key)
	if obj == nil {
		obj = object.ListObject()
		db.set(key, obj)
	}

	ls, ok := obj.List()
	if !ok {
		return 0, ErrWrongTypeOps
	}
	for _, v := range val {
		ls.RPush(v)
	}
	return ls.Length(), nil
}

func (db *DB) LPop(key string) (interface{}, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	ls, ok := obj.List()
	if !ok {
		return nil, ErrWrongTypeOps
	}
	old := ls.LPop()

	if ls.Length() == 0 {
		db.remove(key)
	}
	return old, nil
}

func (db *DB) RPop(key string) (interface{}, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	ls, ok := obj.List()
	if !ok {
		return nil, ErrWrongTypeOps
	}
	old := ls.RPop()

	if ls.Length() == 0 {
		db.remove(key)
	}
	return old, nil
}

func (db *DB) LRange(key string, start, end int) ([]interface{}, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	ls, ok := obj.List()
	if !ok {
		return nil, ErrWrongTypeOps
	}

	startNode := ls.Index(start)
	endNode := ls.Index(end)
	var endNext *list.Node
	if endNode != nil {
		endNext = endNode.Next()
	}
	var vals []interface{}
	for n := startNode; n != nil && n != endNext; n = n.Next() {
		fmt.Println(n.Val())
		vals = append(vals, n.Val())
	}
	return vals, nil
}

func (db *DB) LIndex(key string, index int) (interface{}, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	ls, ok := obj.List()
	if !ok {
		return nil, ErrWrongTypeOps
	}

	n := ls.Index(index)
	if n == nil {
		return nil, nil
	}
	return n.Val(), nil
}

func (db *DB) LSet(key string, index int, newVal interface{}) (interface{}, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	ls, ok := obj.List()
	if !ok {
		return nil, ErrWrongTypeOps
	}

	n := ls.Index(index)
	if n == nil {
		return nil, nil
	}
	old := n.SetVal(newVal)
	return old, nil
}
