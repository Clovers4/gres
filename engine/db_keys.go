package engine

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

func (db *DB) Ttl(key string) int64 {
	return db.ttl(key)
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
