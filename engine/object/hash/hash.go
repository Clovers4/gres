package hash

import (
	"encoding/json"
	"github.com/clovers4/gres/engine/object/plain"
	"github.com/clovers4/gres/util"
	"io"
)

type Hash struct {
	m map[string]interface{} // if needs persistence, val should be the type which plain.Plain support.
}

func New() *Hash {
	return &Hash{
		m: make(map[string]interface{}),
	}
}

func (h *Hash) Set(key string, val interface{}) (interface{}, bool) {
	old, existed := h.m[key]
	h.m[key] = val
	return old, existed
}

func (h *Hash) Delete(key string) (interface{}, bool) {
	old, existed := h.m[key]
	delete(h.m, key)
	return old, existed
}

func (h *Hash) Get(key string) (interface{}, bool) {
	val, existed := h.m[key]
	return val, existed
}

func (h *Hash) Exists(key string) bool {
	_, existed := h.m[key]
	return existed
}

func (h *Hash) Keys() []string {
	var keys []string
	for k := range h.m {
		keys = append(keys, k)
	}
	return keys
}

func (h *Hash) Vals() []interface{} {
	var vals []interface{}
	for _, v := range h.m {
		vals = append(vals, v)
	}
	return vals
}

// 单数是 key, 双数是 val
func (h *Hash) KeyVals() []interface{} {
	var kvs []interface{}
	for k, v := range h.m {
		kvs = append(kvs, k)
		kvs = append(kvs, v)
	}
	return kvs
}

func (h *Hash) Length() int {
	return len(h.m)
}

// Only for test
func (h *Hash) String() string {
	b, err := json.Marshal(h.m)
	if err != nil {
		return "MARSHAL ERROR"
	}
	return string(b)
}

func (h *Hash) Marshal(w io.Writer) error {
	// write total. the total must > 0
	total := h.Length()
	if err := util.Write(w, int64(total)); err != nil {
		return err
	}

	// loop write score and val
	for k, v := range h.m {
		if err := util.Write(w, k); err != nil {
			return err
		}

		// use Plain to marshal
		p := plain.New(v)
		if err := p.Marshal(w); err != nil {
			return err
		}
	}
	return nil
}

func (h *Hash) Unmarshal(r io.Reader) error {
	var total int64
	if err := util.Read(r, &total); err != nil {
		return err
	}

	for i := 0; i < int(total); i++ {
		var key string
		if err := util.Read(r, &key); err != nil {
			return err
		}

		p := plain.New(nil)
		if err := p.Unmarshal(r); err != nil {
			return err
		}

		h.Set(key, p.Val())
	}
	return nil
}
