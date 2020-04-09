package engine

import (
	"github.com/clovers4/gres/engine/object"
	"github.com/clovers4/gres/util"
)

// DbSize cannot get really correct count because of the concurrence.
func (db *DB) DbSize() int {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()

	if db.onSave {
		return db.dataMap.Count() + db.dirtyDataMap.Count()
	}
	return db.dataMap.Count()
}

func (db *DB) Exists(key string) bool {
	return db.get(key) != nil
}

func (db *DB) Ttl(key string) int {
	return int(db.ttl(key))
}

// if return true, the db has old value, otherwise, the db do not has the old kv.
func (db *DB) Del(key ...string) int {
	count := 0
	for _, k := range key {
		if db.remove(k) != nil {
			count++
		}
	}
	return count
}

func (db *DB) Expire(key string, seconds int) bool {
	return db.setExpire(key, seconds)
}

func (db *DB) Type(key string) string {
	obj := db.get(key)
	if obj == nil {
		return "none"
	}
	return obj.Kind().String()
}

func (db *DB) Keys(pattern string) ([]string, error) {
	len := db.DbSize()
	km := make(map[string]bool, len)
	exp := make(map[string]bool, len)
	db.forEachRead(func(key string, val interface{}) {
		if util.Match(pattern, key) {
			if val == object.Expunged {
				exp[key] = true
			} else if !exp[key] {
				km[key] = true
			}
		}
	})

	ks := make([]string, 0, len)
	for k := range km {
		ks = append(ks, k)
	}
	return ks, nil
}
