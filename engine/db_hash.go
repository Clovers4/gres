package engine

import (
	"github.com/clovers4/gres/engine/object"
	"github.com/clovers4/gres/util"
)

// ========
//   Hash
// ========
func (db *DB) HSet(key string, filed string, val interface{}) (int, error) {
	obj := db.get(key)
	if obj == nil {
		obj = object.HashObject()
		db.set(key, obj)
	}

	h, ok := obj.Hash()
	if !ok {
		return 0, ErrWrongTypeOps
	}
	_, existed := h.Set(filed, val)
	if existed {
		return 0, nil
	}
	return 1, nil
}

func (db *DB) HGet(key string, field string) (interface{}, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	h, ok := obj.Hash()
	if !ok {
		return nil, ErrWrongTypeOps
	}
	v, _ := h.Get(field)
	return v, nil
}

func (db *DB) HDel(key string, field ...string) (int, error) {
	obj := db.get(key)
	if obj == nil {
		return 0, nil
	}

	h, ok := obj.Hash()
	if !ok {
		return 0, ErrWrongTypeOps
	}

	count := 0
	for _, f := range field {
		_, existed := h.Delete(f)
		if existed {
			count++
		}
	}
	if h.Length() == 0 {
		db.remove(key)
	}
	return count, nil
}

func (db *DB) HLen(key string) (int, error) {
	obj := db.get(key)
	if obj == nil {
		return 0, nil
	}

	h, ok := obj.Hash()
	if !ok {
		return 0, ErrWrongTypeOps
	}
	len := h.Length()
	return len, nil
}

func (db *DB) HExists(key string, field string) (bool, error) {
	obj := db.get(key)
	if obj == nil {
		return false, nil
	}

	h, ok := obj.Hash()
	if !ok {
		return false, ErrWrongTypeOps
	}
	return h.Exists(field), nil
}

func (db *DB) HKeys(key string) ([]string, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	h, ok := obj.Hash()
	if !ok {
		return nil, ErrWrongTypeOps
	}
	return h.Keys(), nil
}

func (db *DB) HVals(key string) ([]interface{}, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	h, ok := obj.Hash()
	if !ok {
		return nil, ErrWrongTypeOps
	}
	return h.Vals(), nil
}

func (db *DB) HGetAll(key string) ([]interface{}, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	h, ok := obj.Hash()
	if !ok {
		return nil, ErrWrongTypeOps
	}
	return h.KeyVals(), nil
}

func (db *DB) HIncrBy(key string, field string, increment int) (interface{}, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	h, ok := obj.Hash()
	if !ok {
		return nil, ErrWrongTypeOps
	}

	val, ok := h.Get(field)
	if !ok {
		val = 0
	}

	valInt, err := util.IntX2Int(val)
	if err != nil {
		return nil, err
	}

	afterInt, ok := util.Add(valInt, increment)
	if !ok {
		return 0, ErrOutOfRange
	}
	afterVal := util.ShrinkNum(afterInt)
	h.Set(field, afterVal)
	return afterVal, nil
}

// todo:hincrbyfloat
