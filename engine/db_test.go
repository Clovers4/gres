package engine

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/clovers4/gres/engine/object/plain"
	"github.com/clovers4/gres/util"
	"github.com/stretchr/testify/assert"
)

func array2String(vals []interface{}, needSort bool) string {
	var s string
	s += "{"

	var sVal []string
	for _, val := range vals {
		p := plain.New(val)
		sVal = append(sVal, p.String())
	}

	if needSort {
		sort.Strings(sVal)
	}

	for _, val := range sVal {
		s += val + ", "
	}
	if len(s) > 2 {
		s = s[:len(s)-2]
	}
	s += "}"
	return s
}

func TestFile(t *testing.T) {
	var err error
	file, err := os.OpenFile("db_test_file", os.O_CREATE, 0666)
	assert.Nil(t, err)

	w := bufio.NewWriter(file)
	_, err = w.Write(util.StringToBytes("ABC"))
	assert.Nil(t, err)
	w.Flush()
}

func TestDB_Save(t *testing.T) {
	db := NewDB(true)

	db.Set("string-A", "A1")
	db.Set("string-A", "A1-0")
	db.Set("string-B", int32(32))

	db.RPush("list-A", "B")
	db.RPush("list-A", "C")
	db.LPush("list-A", "A")

	db.SAdd("set-A", int32(2))
	db.SAdd("set-A", "SD")

	db.ZAdd("zset-A", 23, "m-A")
	db.ZAdd("zset-A", 12, "m-B")

	db.Save()
	time.Sleep(3 * time.Second)

	db.HSet("hash-A", "f-A", "v-A")
	db.HSet("hash-A", "f-A", "v-A-0")
	db.HSet("hash-A", "f-B", "v-B")
	db.HSet("hash-B", "f-B3", "v-B3")

	fmt.Println(db)
	db.Save()

	newDB := NewDB(true)
	err := newDB.ReadFromFile()
	assert.Nil(t, err)
	fmt.Println(newDB)

	assert.Equal(t, db.String(), newDB.String())
}

func TestDB_ReadFromFile(t *testing.T) {
	newDB := NewDB(true)
	err := newDB.ReadFromFile()
	assert.Nil(t, err)
	fmt.Println(newDB)
}

func TestDB_Plain(t *testing.T) {
	db := NewDB(false)
	var val interface{}
	var err error

	val, err = db.Get("b")
	assert.Nil(t, err)
	assert.Equal(t, nil, val)

	db.Set("b", int64(32))
	val, err = db.Get("b")
	assert.Equal(t, int8(32), val)
	assert.Nil(t, err)

	val, err = db.IncrBy("b", 5)
	assert.Equal(t, int8(37), val)
	assert.Nil(t, err)

	val, err = db.Get("b")
	assert.Equal(t, int8(37), val)
	assert.Nil(t, err)

	val, err = db.Incr("b")
	assert.Equal(t, int8(38), val)
	assert.Nil(t, err)

	val, err = db.Incr("c")
	assert.Equal(t, int8(1), val)
	assert.Nil(t, err)

	val, err = db.DecrBy("b", 4)
	assert.Equal(t, int8(34), val)
	assert.Nil(t, err)

	val, err = db.Decr("b")
	assert.Equal(t, int8(33), val)
	assert.Nil(t, err)

	val, err = db.Get("b")
	assert.Equal(t, int8(33), val)
	assert.Nil(t, err)

	val, err = db.GetSet("b", "ABC")
	assert.Equal(t, int8(33), val)
	assert.Nil(t, err)

	val, err = db.Get("b")
	assert.Equal(t, "ABC", val)
	assert.Nil(t, err)
}

func TestDB_Hash(t *testing.T) {
	db := NewDB(false)
	var val interface{}
	var err error

	val, err = db.HGet("a", "a")
	assert.Nil(t, err)
	assert.Equal(t, nil, val)

	val, err = db.HSet("a", "f-1", "v-1")
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	val, err = db.HGet("a", "f-1")
	assert.Nil(t, err)
	assert.Equal(t, "v-1", val)

	val, err = db.HLen("a")
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	val, err = db.HExists("a", "f-1")
	assert.Nil(t, err)
	assert.Equal(t, true, val)

	val, err = db.HDel("a", "f-1")
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	val, err = db.HLen("a")
	assert.Nil(t, err)
	assert.Equal(t, 0, val)

	val, err = db.HExists("a", "f-1")
	assert.Nil(t, err)
	assert.Equal(t, false, val)

	val, err = db.HDel("a", "f-1")
	assert.Nil(t, err)
	assert.Equal(t, 0, val)

	val, err = db.HSet("a", "f-2", "v-2")
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	val, err = db.HSet("a", "f-3", "v-3")
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	val, err = db.HIncrBy("a", "f-4", 34)
	assert.Nil(t, err)
	assert.Equal(t, int8(34), val)

	keys, err := db.HKeys("a")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(keys))
	fmt.Println(keys)

	vals, err := db.HVals("a")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(vals))
	fmt.Println(vals)

	kvs, err := db.HGetAll("a")
	assert.Nil(t, err)
	assert.Equal(t, 6, len(kvs))
	fmt.Println(kvs)
}

func TestDB_List(t *testing.T) {
	db := NewDB(false)
	var err error
	var val interface{}
	var vals []interface{}

	val, err = db.LPush("ls", "A")
	assert.Nil(t, err)

	val, err = db.RPush("ls", "B")
	assert.Nil(t, err)

	val, err = db.RPush("ls", "C")
	assert.Nil(t, err)

	val, err = db.RPush("ls", 32)
	assert.Nil(t, err)

	vals, err = db.LRange("ls", 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(vals))
	assert.Equal(t, "{B, C}", array2String(vals, false))

	vals, err = db.LRange("ls", 0, -1)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(vals))
	assert.Equal(t, "{A, B, C, 32}", array2String(vals, false))

	val, err = db.LIndex("ls", 1)
	assert.Nil(t, err)
	assert.Equal(t, "B", val)

	val, err = db.LSet("ls", 0, "new-A")
	assert.Nil(t, err)
	assert.Equal(t, "A", val)

	val, err = db.LPop("ls")
	assert.Nil(t, err)
	assert.Equal(t, "new-A", val)

	val, err = db.RPop("ls")
	assert.Nil(t, err)
	assert.Equal(t, 32, val)

	vals, err = db.LRange("ls", 0, -1)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(vals))
	assert.Equal(t, "{B, C}", array2String(vals, false))
}

func TestDB_Set(t *testing.T) {
	db := NewDB(false)
	var err error
	var num int
	var ok bool
	var vals []interface{}

	num, err = db.SAdd("set-1", "A")
	assert.Nil(t, err)
	assert.Equal(t, 1, num)

	db.SAdd("set-1", "B")
	assert.Nil(t, err)

	db.SAdd("set-1", "C")
	assert.Nil(t, err)

	db.SAdd("set-1", 22)
	assert.Nil(t, err)

	db.SAdd("set-1", 5)
	assert.Nil(t, err)

	num, err = db.SCard("set-1")
	assert.Nil(t, err)
	assert.Equal(t, 5, num)

	num, err = db.SCard("set-2")
	assert.Nil(t, err)
	assert.Equal(t, 0, num)

	num, err = db.SRem("set-1", 5)
	assert.Nil(t, err)
	assert.Equal(t, 1, num)

	num, err = db.SRem("set-1", 5)
	assert.Nil(t, err)
	assert.Equal(t, 0, num)

	ok, err = db.SIsMember("set-1", "A")
	assert.Nil(t, err)
	assert.Equal(t, true, ok)

	ok, err = db.SIsMember("set-1", 5)
	assert.Nil(t, err)
	assert.Equal(t, false, ok)

	vals, err = db.SMembers("set-1")
	assert.Nil(t, err)
	assert.Equal(t, "{22, A, B, C}", array2String(vals, true))

	num, err = db.SAdd("set-2", "A")
	assert.Nil(t, err)

	db.SAdd("set-2", "B")
	assert.Nil(t, err)

	db.SAdd("set-2", "D")
	assert.Nil(t, err)

	vals, err = db.SMembers("set-2")
	assert.Nil(t, err)
	assert.Equal(t, "{A, B, D}", array2String(vals, true))

	vals, err = db.SInter("set-1", "set-2")
	assert.Nil(t, err)
	assert.Equal(t, "{22, A, B, C, D}", array2String(vals, true))

	vals, err = db.SUnion("set-1", "set-2")
	assert.Nil(t, err)
	assert.Equal(t, "{A, B}", array2String(vals, true))

	vals, err = db.SDiff("set-1", "set-2")
	assert.Nil(t, err)
	assert.Equal(t, "{22, C}", array2String(vals, true))
}

func TestDB_ZSet(t *testing.T) {
	db := NewDB(false)
	var err error
	var num int
	var rank *int
	var score *float64
	var f float64
	var vals []interface{}

	num, err = db.ZAdd("zs", 1.1, "A")
	assert.Nil(t, err)
	assert.Equal(t, 1, num)

	num, err = db.ZAdd("zs", 2.1, "B")
	assert.Nil(t, err)
	assert.Equal(t, 1, num)

	num, err = db.ZAdd("zs", 3.1, "C")
	assert.Nil(t, err)
	assert.Equal(t, 1, num)

	num, err = db.ZCard("zs")
	assert.Nil(t, err)
	assert.Equal(t, 3, num)

	score, err = db.ZScore("zs", "B")
	assert.Nil(t, err)
	assert.Equal(t, 2.1, *score)

	score, err = db.ZScore("zs", "Unknown")
	assert.Nil(t, err)
	assert.Nil(t, score)

	rank, err = db.ZRank("zs", "B")
	assert.Nil(t, err)
	assert.Equal(t, 1, *rank)

	rank, err = db.ZRank("zs", "Unknown")
	assert.Nil(t, err)
	assert.Nil(t, rank)

	num, err = db.ZRem("zs", "C")
	assert.Nil(t, err)
	assert.Equal(t, 1, num)

	num, err = db.ZRem("zs", "Unknown")
	assert.Nil(t, err)
	assert.Equal(t, 0, num)

	f, err = db.ZIncrBy("zs", 2.5, "A-2")
	assert.Nil(t, err)
	assert.Equal(t, 2.5, f)

	f, err = db.ZIncrBy("zs", 2.5, "A")
	assert.Nil(t, err)
	assert.Equal(t, 3.6, f)

	vals, err = db.ZRange("zs", 0, -1, false)
	assert.Nil(t, err)
	assert.Equal(t, "{B, A-2, A}", array2String(vals, false))

	vals, err = db.ZRange("zs", 0, -1, true)
	assert.Nil(t, err)
	assert.Equal(t, "{B, 2.1, A-2, 2.5, A, 3.6}", array2String(vals, false))

	vals, err = db.ZRange("zs", 1, 1, true)
	assert.Nil(t, err)
	assert.Equal(t, "{A-2, 2.5}", array2String(vals, false))

	vals, err = db.ZRange("zs", 1, 1, false)
	assert.Nil(t, err)
	assert.Equal(t, "{A-2}", array2String(vals, false))
}

func TestDB_Expire(t *testing.T) {
	db := NewDB(true)
	var v interface{}
	var err error
	var ok bool

	err = db.Set("A", "S-1")
	assert.Nil(t, err)

	v, err = db.Get("A")
	assert.Nil(t, err)
	assert.Equal(t, "S-1", v)

	ok = db.Expire("A", 1)
	assert.Equal(t, true, ok)

	ok = db.Expire("C", 2)
	assert.Equal(t, false, ok)

	time.Sleep(3 * time.Second)

	v, err = db.Get("A")
	assert.Nil(t, err)
	//assert.Equal(t, "S-1", v)

	db.Del("A")

	v, err = db.Get("A")
	assert.Nil(t, err)
	assert.Nil(t, v)

	err = db.Set("B", "S-2")
	assert.Nil(t, err)

	v, err = db.Get("B")
	assert.Nil(t, err)
	assert.Equal(t, "S-2", v)

	err = db.Set("C", "S-3")
	assert.Nil(t, err)

	ok = db.Expire("C", 110)
	assert.Equal(t, true, ok)

	fmt.Println("db", db.String())

	fmt.Println("ttl:", db.Ttl("A"))
	fmt.Println("ttl:", db.Ttl("B"))
	fmt.Println("ttl:", db.Ttl("C"))

	time.Sleep(2 * time.Second)
}

func TestDB_Expire2(t *testing.T) {
	db := NewDB(true)
	var val interface{}
	var err error

	fmt.Println(db.String())

	val, err = db.Get("A")
	fmt.Println(val, err)
	val, err = db.Get("B")
	fmt.Println(val, err)
	val, err = db.Get("C")
	fmt.Println(val, err)
	fmt.Println()

	fmt.Println("ttl:", db.Ttl("A"))
	fmt.Println("ttl:", db.Ttl("B"))
	fmt.Println("ttl:", db.Ttl("C"))

	time.Sleep(2 * time.Second)
}
