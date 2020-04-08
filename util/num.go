package util

import (
	"github.com/clovers4/gres/errs"
	"math"
	"strconv"
)

func Add(a, b int) (int, bool) {
	c := a + b
	if (c > a) == (b > 0) {
		return c, true
	}
	return c, false
}

func IntX2Int(num interface{}) (int, error) {
	switch v := num.(type) {
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case int:
		return int(v), nil
	default:
		return 0, errs.ErrIsNotInt
	}
}

func ShrinkNum(num int) interface{} {
	if math.MinInt8 <= num && num <= math.MaxInt8 {
		return int8(num)
	} else if math.MinInt16 <= num && num <= math.MaxInt16 {
		return int16(num)
	} else if math.MinInt32 <= num && num <= math.MaxInt32 {
		return int32(num)
	} else {
		return int64(num)
	}
}

func String2Num(s string) (num interface{}, success bool) {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i, true
	}

	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, true
	}
	return nil, false
}

func String2Int(s string) (num int, err error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return int(i), err
}

func String2Float(s string) (num float64, err error) {
	return strconv.ParseFloat(s, 64)
}
