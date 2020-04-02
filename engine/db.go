package engine

import (
	"github.com/clovers4/gres/container/cmap"
	"path/filepath"
	"sort"
	"sync"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
)

const FilenameFormat = "gres_%v"
const FilenameRegex = "gres_*"
const Version int32 = 1

var expunged = newObject(ObjPlain, "expunged")

type DB struct {
	persist  bool // 是否要持久化
	filename string

	cleanMap *cmap.CMap // 正常情况下, 往该 map 中进行存取

	onSave    bool // 持久化中
	dirtyLock sync.RWMutex
	dirtyMap  *cmap.CMap // 持久化中, 新数据存入该 map
}

func NewDB(persist bool) *DB {
	if persist {

	}

	return &DB{
		persist: persist,

		cleanMap: cmap.New(),
	}
}

func (db *DB) set(key string, obj *Object) *Object {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()

	targetMap := db.cleanMap
	if db.onSave {
		targetMap = db.dirtyMap
	}

	if oldValue, existed := targetMap.Set(key, obj); existed {
		return oldValue.(*Object)
	}
	return nil
}

func (db *DB) remove(key string) *Object {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()

	// 持久化中
	if db.onSave {
		// 设置 expunged 作为空标志位
		if oldValue, existed := db.dirtyMap.Set(key, expunged); existed {
			return oldValue.(*Object)
		}

		// 从 cleanMap 读取原始值
		if oldValue, existed := db.cleanMap.Get(key); existed {
			return oldValue.(*Object)
		}
		return nil
	}

	// 非持久化中
	if oldValue, existed := db.cleanMap.Remove(key); existed {
		return oldValue.(*Object)
	}
	return nil
}

func (db *DB) get(key string) *Object {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()

	// 持久化中
	if db.onSave {
		// 先从 dirtyMap 读, 若为 expunged, 则认为已删除
		if oldValue, existed := db.dirtyMap.Get(key); existed {
			if oldValue == expunged {
				return nil
			}
			return oldValue.(*Object)
		}
	}

	// 非持久化中默认读 cleanMap; 持久化中, 若 dirtyMap 无数据, 则从 cleanMap 读取
	v, ok := db.cleanMap.Get(key)
	if !ok {
		return nil
	}
	return v.(*Object)
}

func (db *DB) loadFile() error {
	// 爬文件(默认已经是排序好的)
	files, err := filepath.Glob("my_*")
	if err != nil {
		return err
	}
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	// 校验文件
	return nil
}

// 保证即使持久化过程中断电, 本地文件保存的数据仍具有一致性,
//func (db *DB) Save() error {
//	db.startSave()
//	defer db.endSave()
//
//	newFilename := fmt.Sprintf(FilenameFormat, time.Now().Unix())
//	newFile, err := os.OpenFile(newFilename, os.O_CREATE, 0666)
//	if err != nil {
//		return err
//	}
//
//	// write cleanMap to file
//	count, kvCh := db.cleanMap.Snapshot()
//	for i := 0; i < count; i++ {
//		kv := <-kvCh
//		k, v := kv.Key, kv.Value
//
//	}
//
//	// flush dirtyMap to cleanMap: 需要放在持久化完成之后
//	dirtyCount, dirtyKvCh := db.dirtyMap.Snapshot()
//	for i := 0; i < dirtyCount; i++ {
//		kv := <-dirtyKvCh
//		k, v := kv.Key, kv.Value
//		if v == expunged {
//			db.cleanMap.Remove(k)
//		} else {
//			db.cleanMap.Set(k, v)
//		}
//	}
//
//	return nil
//}

func (db *DB) startSave() {
	db.dirtyLock.Lock()
	db.onSave = true
	db.dirtyMap = cmap.New()
}

func (db *DB) endSave() {
	db.dirtyLock.Unlock()
	db.onSave = false
	db.dirtyMap = nil
}

func (db *DB) CheckKind(key string, kind ObjKind) bool {
	obj := db.get(key)
	if obj == nil {
		return true
	}
	if obj.Kind == kind {
		return true
	}
	return false
}

// ========
// Commands
// ========

//func (db *DB) Set(key string, val interface{}) *Object {
//	obj := PlainObject(val)
//	db.set(cli.args[1], obj)
//
//	return obj
//}
