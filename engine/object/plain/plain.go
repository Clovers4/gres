package plain

import (
	"fmt"
	"github.com/clovers4/gres/util"
	"io"
	"reflect"
)

// only support int8/int16/int32/int64 uint8/uint16/uint32/uint64 float32/float64 string
// NOT support int/uint
type Plain struct {
	val interface{}
}

func New(val interface{}) *Plain {
	return &Plain{
		val: val,
	}
}

func (p *Plain) Val() interface{} {
	return p.val
}

// Only for test
func (p *Plain) String() string {
	return fmt.Sprintf("%v", p.val)
}

func (p *Plain) Marshal(w io.Writer) error {
	// 经常使用 int 忘了转成 int64, 这里做一层防御
	if v, ok := p.val.(int); ok {
		p.val = int64(v)
	}

	kind := uint8(reflect.TypeOf(p.val).Kind())
	if err := util.Write(w, kind); err != nil {
		return err
	}
	if err := util.Write(w, p.val); err != nil {
		return err
	}
	return nil
}

func (p *Plain) Unmarshal(r io.Reader) error {
	var kind uint8
	if err := util.Read(r, &kind); err != nil {
		return err
	}

	switch reflect.Kind(kind) {
	case reflect.Bool:
		var v bool
		if err := util.Read(r, &v); err != nil {
			return err
		}
		p.val = v
	case reflect.Int8:
		var v int8
		if err := util.Read(r, &v); err != nil {
			return err
		}
		p.val = v
	case reflect.Int16:
		var v int16
		if err := util.Read(r, &v); err != nil {
			return err
		}
		p.val = v
	case reflect.Int32:
		var v int32
		if err := util.Read(r, &v); err != nil {
			return err
		}
		p.val = v
	case reflect.Int64:
		var v int64
		if err := util.Read(r, &v); err != nil {
			return err
		}
		p.val = v
	case reflect.Uint8:
		var v uint8
		if err := util.Read(r, &v); err != nil {
			return err
		}
		p.val = v
	case reflect.Uint16:
		var v uint16
		if err := util.Read(r, &v); err != nil {
			return err
		}
		p.val = v
	case reflect.Uint32:
		var v uint32
		if err := util.Read(r, &v); err != nil {
			return err
		}
		p.val = v
	case reflect.Uint64:
		var v uint64
		if err := util.Read(r, &v); err != nil {
			return err
		}
		p.val = v
	case reflect.Float32:
		var v float32
		if err := util.Read(r, &v); err != nil {
			return err
		}
		p.val = v
	case reflect.Float64:
		var v float64
		if err := util.Read(r, &v); err != nil {
			return err
		}
		p.val = v
	case reflect.String:
		var v string
		if err := util.Read(r, &v); err != nil {
			return err
		}
		p.val = v
	default:
		return fmt.Errorf("unsupported plain type [%v]", kind)
	}
	return nil
}
