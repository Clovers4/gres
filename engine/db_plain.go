package engine

import (
	"errors"
	"github.com/clovers4/gres/engine/object"
	"github.com/clovers4/gres/util"
)

var (
	ErrOutOfRange = errors.New("the value will be out of range")
)

// ==============================
//            Plain(String)
// ==============================
func (db *DB) Set(key string, val interface{}) error {
	// 若是数字, 将数字调节到合适的大小以节省内存
	num, err := util.IntX2Int(val)
	if err == nil {
		val = util.ShrinkNum(num)
	}

	obj := object.PlainObject(val)
	db.set(key, obj)
	return nil
}

func (db *DB) Get(key string) (val interface{}, err error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}
	p, ok := obj.Plain()
	if !ok {
		return nil, ErrWrongTypeOps
	}
	return p.Val(), nil
}

func (db *DB) GetSet(key string, val interface{}) (oldVal interface{}, err error) {
	oldVal, err = db.Get(key)
	if err != nil {
		return nil, err
	}
	err = db.Set(key, val)
	return oldVal, err
}

func (db *DB) Incr(key string) (interface{}, error) {
	return db.IncrBy(key, 1)
}

// we think num is always int, and do not use uint.
// 返回结果为计算后的值
func (db *DB) IncrBy(key string, num int) (interface{}, error) {
	obj := db.get(key)
	if obj == nil {
		obj = object.PlainObject(int8(0))
		db.set(key, obj)
	}

	p, ok := obj.Plain()
	if !ok {
		return 0, ErrWrongTypeOps
	}

	oldVal := p.Val()
	old, err := util.IntX2Int(oldVal)
	if err != nil {
		return 0, err
	}

	afterInt, ok := util.Add(old, num)
	if !ok {
		return 0, ErrOutOfRange
	}

	afterVal := util.ShrinkNum(afterInt)
	p.SetVal(afterVal)
	return afterVal, nil
}

func (db *DB) DecrBy(key string, num int) (interface{}, error) {
	return db.IncrBy(key, -num)
}

func (db *DB) Decr(key string) (interface{}, error) {
	return db.IncrBy(key, -1)
}
