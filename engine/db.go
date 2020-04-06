package engine

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/clovers4/gres/engine/cmap"
	"github.com/clovers4/gres/engine/object"
	"github.com/clovers4/gres/util"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
)

var (
	// file
	ErrCRCNotEqual        = errors.New("CRC is not equal to the expect, the file may be broken")
	ErrUnexpectHeader     = errors.New("the header is unexpect")
	ErrUnsupportedVersion = errors.New("the version is unsupported")

	// commands
	ErrWrongTypeOps = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
)

const (
	FilenamePrefix = "gres_"
	FilenameFormat = FilenamePrefix + "%v.db"
	FilenameRegex  = FilenamePrefix + "*"
	GRES           = "GRES"
	DBVersion      = "0.0.1"
)

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

func (db *DB) set(key string, obj *object.Object) *object.Object {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()

	targetMap := db.cleanMap
	if db.onSave {
		targetMap = db.dirtyMap
	}

	if oldValue, existed := targetMap.Set(key, obj); existed {
		return oldValue.(*object.Object)
	}
	return nil
}

func (db *DB) remove(key string) *object.Object {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()

	// 持久化中
	if db.onSave {
		// 设置 Expunged 作为空标志位
		if oldValue, existed := db.dirtyMap.Set(key, object.Expunged); existed {
			return oldValue.(*object.Object)
		}

		// 从 cleanMap 读取原始值
		if oldValue, existed := db.cleanMap.Get(key); existed {
			return oldValue.(*object.Object)
		}
		return nil
	}

	// 非持久化中
	if oldValue, existed := db.cleanMap.Remove(key); existed {
		return oldValue.(*object.Object)
	}
	return nil
}

func (db *DB) get(key string) *object.Object {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()

	// 持久化中
	if db.onSave {
		// 先从 dirtyMap 读, 若为 Expunged, 则认为已删除
		if oldValue, existed := db.dirtyMap.Get(key); existed {
			if oldValue == object.Expunged {
				return nil
			}
			return oldValue.(*object.Object)
		}
	}

	// 非持久化中默认读 cleanMap; 持久化中, 若 dirtyMap 无数据, 则从 cleanMap 读取
	v, ok := db.cleanMap.Get(key)
	if !ok {
		return nil
	}
	return v.(*object.Object)
}

// 保证即使持久化过程中断电, 本地文件保存的数据仍具有一致性,
func (db *DB) Save() error {
	db.dirtyLock.Lock()
	db.onSave = true
	db.dirtyMap = cmap.New()

	// open file
	newFilename := fmt.Sprintf(FilenameFormat, time.Now().Unix())
	newFile, err := os.OpenFile(newFilename, os.O_CREATE, 0666)
	if err != nil {
		db.onSave = false
		db.dirtyMap = nil
		db.dirtyLock.Unlock()
		return err
	}
	defer newFile.Close()
	db.dirtyLock.Unlock()

	// save data to new file
	err = db.save(newFile)

	// end save
	db.dirtyLock.Lock()
	db.cleanMap.AddCMap(db.dirtyMap) // flush dirtyMap to cleanMap: 需要放在持久化完成之后. 此时, db 的 set/get 无法使用，直到完成
	db.onSave = false
	db.dirtyMap = nil

	// 删除老文件
	if err = os.Remove(db.filename); err != nil {
		// todo:log
	}
	db.filename = newFilename
	db.dirtyLock.Unlock()
	return err
}

func (db *DB) save(file *os.File) error {
	var err error

	// write cleanMap to file. Even if failed, needs to write dirtyMap to cleanMap
	w := NewCRCWriter(bufio.NewWriter(file))

	// write constant "GRES" and DB_VERSION
	if err := util.Write(w, GRES); err != nil {
		return err
	}

	// write DB_VERSION
	if err := util.Write(w, DBVersion); err != nil {
		return err
	}

	// write data
	if err = db.cleanMap.Marshal(w); err != nil {
		return err
	}

	// write crc
	if err = w.WriteCRC(); err != nil {
		return err
	}

	if err = w.Flush(); err != nil {
		return err
	}

	return nil
}

func (db *DB) readFromFile(filename string) error {
	var err error
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	r := NewCRCReader(bufio.NewReader(file))

	// read header
	var gresFlag string
	if err := util.Read(r, &gresFlag); err != nil {
		return err
	} else if gresFlag != GRES {
		return ErrUnexpectHeader
	}

	// read version
	var dbVersion string
	if err := util.Read(r, &dbVersion); err != nil {
		return err
	} else if dbVersion != DBVersion {
		return ErrUnsupportedVersion
	}

	// read cleanMap
	if err = db.cleanMap.Unmarshal(r); err != nil {
		return err
	}

	// read crc and check whether is equal to the expect
	if equal, err := r.IsCRCEqual(); err != nil {
		return err
	} else if !equal {
		return ErrCRCNotEqual
	}
	return nil
}

func (db *DB) ReadFromFile() error {
	var err error
	filenames, err := filepath.Glob(FilenameRegex)
	if err != nil {
		return err
	}

	// 如果持久化过程中断电 or 其他极端情况, 可能出现多个 .db 文件
	var stamps []int
	for _, filename := range filenames {
		if len(filename) < len(FilenamePrefix) {
			continue
		}

		num, err := strconv.Atoi(filename[len(FilenamePrefix) : len(filename)-3])
		if err != nil {
			//todo:
			return err
		}
		stamps = append(stamps, num)
	}

	// 从最新的开始读取; 不做删除操作, 以便保留手动修复文件的可能性
	sort.Sort(sort.Reverse(sort.IntSlice(stamps)))
	for i, stamp := range stamps {
		filename := fmt.Sprintf(FilenameFormat, stamp)
		if err = db.readFromFile(filename); err != nil {
			//todo: log
			if i < len(stamps)-1 {
				continue
			}
			return err
		} else {
			db.filename = filename
			break
		}
	}

	return nil
}

// Only for test
func (db *DB) String() string {
	return db.cleanMap.String()
}

// ===============================================
//                   Commands
// ===============================================

// DbSize cannot get really correct count because of the concurrence.
func (db *DB) DbSize() int {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()

	if db.onSave {
		return db.cleanMap.Count() + db.dirtyMap.Count()
	}
	return db.cleanMap.Count()
}

func (db *DB) Exists(key string) bool {
	return db.get(key) != nil
}

// if return true, the db has old value, otherwise, the db do not has the old kv.
func (db *DB) Del(key string) bool {
	return db.remove(key) != nil
}

func (db *DB) Expire(key string) {
	// todo
}

func (db *DB) Type(key string) string {
	obj := db.get(key)
	if obj == nil {
		return "none"
	}
	return obj.Kind().String()
}
