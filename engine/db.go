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
	"github.com/clovers4/gres/zset"
	"github.com/gofrs/flock"
	"go.uber.org/zap"
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
	persist     bool          // 是否要持久化
	persistTime time.Duration // [persist策略] 每隔多久执行一次持久化

	doExpireTime       time.Duration // [expire策略] 每隔多久执行一次 doExpire
	doExpireMinNum     int           // [expire策略] 执行 expire 最少个数
	doExpireMinPercent float64       // [expire策略] 执行 expire 最小百分比

	filename string
	fileLock *flock.Flock

	dataMap    *cmap.CMap // 正常情况下, 往该 map 中进行存取
	expireList *zset.ZSet // 实现过期功能. k=key(string),v=time(unixtime-int64)

	onSave          bool // 持久化中
	dirtyLock       sync.RWMutex
	dirtyDataMap    *cmap.CMap // 持久化中, 新数据存入该 map
	dirtyExpireList *zset.ZSet // 持久化中, 新数据存入该 map

	log *zap.Logger
}

type dbOption func(*DB)

func PersistOption(persist bool) dbOption {
	return func(db *DB) {
		db.persist = persist
	}
}

func LogOption(log *zap.Logger) dbOption {
	return func(db *DB) {
		db.log = log
	}
}

func NewDB(ops ...dbOption) *DB {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	db := &DB{
		persist:     false,
		persistTime: 1 * time.Second, // default

		doExpireTime:       1 * time.Second, // default
		doExpireMinNum:     100,             // default
		doExpireMinPercent: 0.20,            // default

		dataMap:    cmap.New(),
		expireList: zset.New(),
		log:        log,
	}
	for _, op := range ops {
		op(db)
	}

	if db.persist {
		if err := db.keepOneProcess(); err != nil {
			panic(err)
		}
		if err := db.ReadFromFile(); err != nil {
			panic(err)
		}
		go db.SaveBackground()
	}

	go db.DoExpireBackground()

	return db
}

func (db *DB) set(key string, obj *object.Object) *object.Object {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()

	return db.setLocked(key, obj)
}

func (db *DB) setLocked(key string, obj *object.Object) *object.Object {
	targetMap := db.dataMap
	if db.onSave {
		targetMap = db.dirtyDataMap
	}

	if oldValue, existed := targetMap.Set(key, obj); existed {
		return oldValue.(*object.Object)
	}
	return nil
}

func (db *DB) remove(key string) *object.Object {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()
	return db.removeLocked(key)
}

func (db *DB) removeLocked(key string) *object.Object {
	// 持久化中
	if db.onSave {
		// 设置 Expunged 作为空标志位
		if oldValue, existed := db.dirtyDataMap.Set(key, object.Expunged); existed {
			return oldValue.(*object.Object)
		}

		// 从 dataMap 读取原始值
		if oldValue, existed := db.dataMap.Get(key); existed {
			return oldValue.(*object.Object)
		}
		return nil
	}

	// 非持久化中
	if oldValue, existed := db.dataMap.Remove(key); existed {
		return oldValue.(*object.Object)
	}
	return nil
}

func (db *DB) get(key string) *object.Object {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()

	return db.getLocked(key)
}

func (db *DB) getLocked(key string) *object.Object {
	// 持久化中
	if db.onSave {
		// use the side effect, if key is expired, then kv will be nil
		db.ttlLocked(key)

		// 先从 dirtyDataMap 读, 若为 Expunged, 则认为已删除
		if oldValue, existed := db.dirtyDataMap.Get(key); existed {
			if oldValue == object.Expunged {
				return nil
			}
			return oldValue.(*object.Object)
		}
	}

	// 非持久化中默认读 dataMap; 持久化中, 若 dirtyDataMap 无数据, 则从 dataMap 读取
	v, ok := db.dataMap.Get(key)
	if !ok {
		return nil
	}
	return v.(*object.Object)
}

// heavy ops, only for test, keys command....
func (db *DB) forEachRead(fn func(key string, val interface{})) {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()
	db.forEachReadLocked(fn)
}

func (db *DB) forEachReadLocked(fn func(key string, val interface{})) {
	// 持久化中
	if db.onSave {
		db.dirtyDataMap.ForEachRead(fn)
	}

	// 非持久化中默认读 dataMap; 持久化中, 若 dirtyDataMap 无数据, 则从 dataMap 读取
	db.dataMap.ForEachRead(fn)
}

func (db *DB) setExpire(key string, seconds int) bool {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()
	return db.setExpireLocked(key, seconds)
}

func (db *DB) setExpireLocked(key string, seconds int) bool {
	// 若传入 0 秒, 则直接删除
	if seconds == 0 {
		if db.removeLocked(key) == nil {
			return false
		}
		return true
	}

	if db.getLocked(key) == nil {
		return false
	}

	targetList := db.expireList
	if db.onSave {
		targetList = db.dirtyExpireList
	}

	endTime := time.Now().Add(time.Duration(seconds) * time.Second).Unix()
	return targetList.Add(endTime, key)
}

func (db *DB) removeExpire(key string) bool {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()
	return db.removeExpireLocked(key)
}

func (db *DB) removeExpireLocked(key string) bool {
	// 持久化中
	if db.onSave {
		// 设置 -1 作为空标志位
		if notExist := db.dirtyExpireList.Add(-1, key); !notExist {
			return true
		}
		return false
	}

	// 非持久化中
	if _, existed := db.expireList.Delete(key); existed {
		return true
	}
	return false
}

func (db *DB) ttl(key string) int64 {
	db.dirtyLock.RLock()
	defer db.dirtyLock.RUnlock()
	return db.ttlLocked(key)
}

func (db *DB) ttlLocked(key string) int64 {
	// 持久化中
	if db.onSave {
		// 先从 dirtyExpireList 读, 若为 Expunged, 则认为 expire 记录已删除
		if t, existed := db.dirtyExpireList.Get(key); existed {
			if t == -1 {
				// 再查看 dataMap 是否有数据
				if db.getLocked(key) == nil {
					return -2 // key 不存在但无 expire 记录
				} else {
					return -1 // key 存在但无 expire 记录
				}
			}

			now := time.Now().Unix()
			if now >= t {
				// 说明已过期
				db.removeExpireLocked(key)
				db.removeLocked(key)
				return -2
			}
			return t - now
		}
	}

	// 非持久化中默认读 dataMap; 持久化中, 若 dirtyDataMap 无数据, 则从 dataMap 读取
	t, ok := db.expireList.Get(key)
	if !ok {
		if db.getLocked(key) == nil {
			return -2 // key 不存在但无 expire 记录
		} else {
			return -1 // key 存在但无 expire 记录
		}
	}
	now := time.Now().Unix()
	if now >= t {
		// 说明已过期
		db.removeExpireLocked(key)
		db.removeLocked(key)
		return -2
	}
	return t - now
}

func (db *DB) DoExpireBackground() {
	t := time.NewTicker(db.doExpireTime)
	for {
		<-t.C
		db.doExpire()
	}
}

func (db *DB) doExpire() {
	db.dirtyLock.RLock()

	now := time.Now().Unix()
	target := db.expireList

	// 持久化中
	if db.onSave {
		target = db.dirtyExpireList
	}

	target.RLock()
	n := target.GetNodeLeastScore(0) // ignore -1
	loop := 0
	minNum := int(float64(target.Length()) * db.doExpireMinPercent)
	if db.doExpireMinNum < minNum {
		minNum = db.doExpireMinNum
	}
	needDels := make([]string, 0, minNum)

	for {
		if n == nil {
			break
		}
		if loop > minNum {
			break
		}

		// 已过期
		if now >= n.Score() {
			needDels = append(needDels, n.Val())
		} else {
			break
		}

		n = n.Next()
	}
	target.RUnlock()
	db.dirtyLock.RUnlock()

	for _, key := range needDels {
		db.dirtyLock.RLock()
		db.removeExpire(key)
		db.remove(key)
		db.dirtyLock.RUnlock()
	}

	db.log.Debug("[DB doExpire] finished", zap.String("now", time.Now().String()))
}

func (db *DB) SaveBackground() {
	t := time.NewTicker(db.persistTime)
	for {
		<-t.C
		if err := db.Save(); err != nil {
			db.log.Error("[DB SaveBackground] Save", zap.String("err", err.Error()))
		} else {
			db.log.Debug("[DB SaveBackground] Save success")
		}
	}
}

// 保证即使持久化过程中断电, 本地文件保存的数据仍具有一致性,
func (db *DB) Save() error {
	db.dirtyLock.Lock()
	db.onSave = true
	db.dirtyDataMap = cmap.New()
	db.dirtyExpireList = zset.New()

	// open file
	newFilename := fmt.Sprintf(FilenameFormat, time.Now().Unix())
	newFile, err := os.OpenFile(newFilename, os.O_CREATE, 0666)
	if err != nil {
		db.onSave = false
		db.dirtyDataMap = nil
		db.dirtyLock.Unlock()
		return err
	}
	defer newFile.Close()
	db.dirtyLock.Unlock()

	// save data to new file
	err = db.save(newFile)

	// end save
	db.dirtyLock.Lock()
	db.dataMap.AddCMap(db.dirtyDataMap)       // flush dirtyDataMap to dataMap: 需要放在持久化完成之后. 此时, db 的 set/get 无法使用，直到完成
	db.expireList.AddZSet(db.dirtyExpireList) // flush dirtyExpireList to expireList: 需要放在持久化完成之后. 此时, db 的 expire/... 无法使用，直到完成
	db.onSave = false
	db.dirtyDataMap = nil
	db.dirtyExpireList = nil

	// 删除老文件
	if db.filename != "" {
		if err := os.Remove(db.filename); err != nil {
			db.log.Error("[DB Save] Remove", zap.String("err", err.Error()))
		}
	}

	db.filename = newFilename
	db.dirtyLock.Unlock()
	return err
}

func (db *DB) save(file *os.File) error {
	var err error

	// write dataMap to file. Even if failed, needs to write dirtyDataMap to dataMap
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
	if err = db.dataMap.Marshal(w); err != nil {
		return err
	}

	// write expire
	if err = db.expireList.Marshal(w); err != nil {
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

	// read dataMap
	if err = db.dataMap.Unmarshal(r); err != nil {
		return err
	}

	// write expire
	if err = db.expireList.Unmarshal(r); err != nil {
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
	for i, filename := range filenames {
		if len(filename) < len(FilenamePrefix) {
			continue
		}

		num, err := strconv.Atoi(filename[len(FilenamePrefix) : len(filename)-3])
		if err != nil {
			db.log.Error("[DB ReadFromFile] read stamp", zap.String("err", err.Error()))
			if i == len(filenames)-1 && len(stamps) == 0 {
				return err
			}
		}
		stamps = append(stamps, num)
	}

	// 从最新的开始读取; 不做删除操作, 以便保留手动修复文件的可能性
	sort.Sort(sort.Reverse(sort.IntSlice(stamps)))
	for i, stamp := range stamps {
		filename := fmt.Sprintf(FilenameFormat, stamp)
		if err = db.readFromFile(filename); err != nil {
			db.log.Error("[DB ReadFromFile] readFromFile", zap.String("err", err.Error()))
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

func (db *DB) Close() error {
	var err error
	db.doExpire()
	err = db.Save()

	if db.persist {
		db.endKeepOneProcess()
	}

	return err
}

func (db *DB) keepOneProcess() error {
	db.fileLock = flock.New("GRES_LOCK")
	ok, err := db.fileLock.TryLock()
	if !ok {
		return fmt.Errorf("GRES is already boost")
	}
	return err
}

func (db *DB) endKeepOneProcess() error {
	return db.fileLock.Unlock()
}

// Only for test
func (db *DB) String() string {
	return db.dataMap.String()
}
