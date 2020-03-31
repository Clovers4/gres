package engine

import (
	"errors"
	"reflect"
	"strconv"
)

var ErrNum2StringType = errors.New("unsupported type in num2string")

func String2Num(s string) (num interface{}, err error) {
	num, err = strconv.ParseInt(s, 10, 64)
	if err == nil {
		return num, nil
	}

	num, err = strconv.ParseFloat(s, 64)
	if err == nil {
		return num, nil
	}
	return nil, err
}

func Num2String(num interface{}) (string, error) {
	switch reflect.TypeOf(num).Kind() {
	case reflect.Int64:
		numInt64 := num.(int64)
		s := strconv.FormatInt(numInt64, 10)
		return s, nil
	case reflect.Float64:
		numFloat64 := num.(float64)
		s := strconv.FormatFloat(numFloat64, 'f', -1, 64)
		return s, nil
	default:
		return "", ErrNum2StringType
	}
}
