package engine

import (
	"github.com/clovers4/gres/engine/object"
	"github.com/clovers4/gres/engine/object/set"
)

// =========
//    Set
// =========
func (db *DB) SAdd(key string, val ...interface{}) (int, error) {
	obj := db.get(key)
	if obj == nil {
		obj = object.SetObject()
		db.set(key, obj)
	}

	set, ok := obj.Set()
	if !ok {
		return 0, ErrWrongTypeOps
	}
	count := 0
	for _, v := range val {
		if set.Add(v) {
			count++
		}
	}
	return count, nil
}

func (db *DB) SRem(key string, val ...interface{}) (int, error) {
	obj := db.get(key)
	if obj == nil {
		return 0, nil
	}

	set, ok := obj.Set()
	if !ok {
		return 0, ErrWrongTypeOps
	}

	count := 0
	for _, v := range val {
		_, existed := set.Delete(v)
		if existed {
			count++
		}
	}
	if set.Length() == 0 {
		db.remove(key)
	}
	return count, nil
}

func (db *DB) SCard(key string) (int, error) {
	obj := db.get(key)
	if obj == nil {
		return 0, nil
	}

	set, ok := obj.Set()
	if !ok {
		return 0, ErrWrongTypeOps
	}
	len := set.Length()
	return len, nil
}

func (db *DB) SIsMember(key string, val interface{}) (bool, error) {
	obj := db.get(key)
	if obj == nil {
		return false, nil
	}

	set, ok := obj.Set()
	if !ok {
		return false, ErrWrongTypeOps
	}
	existed := set.Exists(val)
	return existed, nil
}

func (db *DB) SMembers(key string) ([]interface{}, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	set, ok := obj.Set()
	if !ok {
		return nil, ErrWrongTypeOps
	}
	return set.Vals(), nil
}

func (db *DB) SInter(keys ...string) ([]interface{}, error) {
	inter := set.New()
	for i, key := range keys {
		obj := db.get(key)
		if obj == nil {
			return nil, nil
		}
		set, ok := obj.Set()
		if !ok {
			return nil, ErrWrongTypeOps
		}
		if i == 0 {
			inter = set
		} else {
			inter = inter.Inter(set)
		}
	}
	return inter.Vals(), nil
}

func (db *DB) SUnion(keys ...string) ([]interface{}, error) {
	union := set.New()
	for i, key := range keys {
		obj := db.get(key)
		if obj == nil {
			return nil, nil
		}
		set, ok := obj.Set()
		if !ok {
			return nil, ErrWrongTypeOps
		}
		if i == 0 {
			union = set
		} else {
			union = union.Union(set)
		}
	}
	return union.Vals(), nil
}

func (db *DB) SDiff(keys ...string) ([]interface{}, error) {
	diff := set.New()
	for i, key := range keys {
		obj := db.get(key)
		if obj == nil {
			return nil, nil
		}
		set, ok := obj.Set()
		if !ok {
			return nil, ErrWrongTypeOps
		}
		if i == 0 {
			diff = set
		} else {
			diff = diff.Diff(set)
		}
	}
	return diff.Vals(), nil
}
