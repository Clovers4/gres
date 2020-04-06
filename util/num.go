package util

import (
	"errors"
	"math"
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
		return 0, errors.New("the value is not an integer")
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
