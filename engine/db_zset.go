package engine

import (
	"github.com/clovers4/gres/engine/object"
)

// ========
//   ZSet
// ========
func (db *DB) ZAdd(key string, score float64, member string) (int, error) {
	obj := db.get(key)
	if obj == nil {
		obj = object.ZSetObject()
		db.set(key, obj)
	}

	zs, ok := obj.ZSet()
	if !ok {
		return 0, ErrWrongTypeOps
	}
	if zs.Add(score, member) {
		return 1, nil
	}
	return 0, nil
}

func (db *DB) ZCard(key string) (int, error) {
	obj := db.get(key)
	if obj == nil {
		return 0, nil
	}

	zs, ok := obj.ZSet()
	if !ok {
		return 0, ErrWrongTypeOps
	}
	return zs.Length(), nil
}

func (db *DB) ZScore(key, member string) (*float64, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	zs, ok := obj.ZSet()
	if !ok {
		return nil, ErrWrongTypeOps
	}
	score, existed := zs.Get(member)
	if !existed {
		return nil, nil
	}
	return &score, nil
}

func (db *DB) ZRank(key, member string) (*int, error) {
	obj := db.get(key)
	if obj == nil {
		return nil, nil
	}

	zs, ok := obj.ZSet()
	if !ok {
		return nil, ErrWrongTypeOps
	}

	rank, existed := zs.GetRankByMember(member)
	if !existed {
		return nil, nil
	}
	return &rank, nil
}

func (db *DB) ZRem(key string, member ...string) (int, error) {
	obj := db.get(key)
	if obj == nil {
		return 0, nil
	}

	zs, ok := obj.ZSet()
	if !ok {
		return 0, ErrWrongTypeOps
	}

	count := 0
	for _, m := range member {
		if _, existed := zs.Delete(m); existed {
			count++
		}
	}

	if zs.Length() == 0 {
		db.remove(key)
	}
	return count, nil
}

// todo: panic
func (db *DB) ZIncrBy(key string, increment float64, member string) (float64, error) {
	obj := db.get(key)
	if obj == nil {
		obj = object.ZSetObject()
		db.set(key, obj)
	}

	zs, ok := obj.ZSet()
	if !ok {
		return 0, ErrWrongTypeOps
	}

	score, _ := zs.Get(member) // if not existed, the score is 0
	zs.Add(score+increment, member)
	return score + increment, nil
}

func (db *DB) ZRange(key string, start, end int, withScore bool) ([]interface{}, error) {
	obj := db.get(key)
	if obj == nil {
		obj = object.ZSetObject()
		db.set(key, obj)
	}

	zs, ok := obj.ZSet()
	if !ok {
		return nil, ErrWrongTypeOps
	}

	startNode := zs.GetNodeByRank(start)
	endNode := zs.GetNodeByRank(end)
	var vals []interface{}
	for n := startNode; n != nil && n != endNode.Next(); n = n.Next() {
		vals = append(vals, n.Val())
		if withScore {
			vals = append(vals, n.Score())
		}
	}
	return vals, nil
}

// todo: zcount zremrangebyrank zremrangebystore
