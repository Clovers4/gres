package set

import (
	"fmt"
	"github.com/clovers4/gres/engine/object/plain"
	"github.com/clovers4/gres/util"
	"io"
	"sort"
)

//todo

// only support int8/int16/int32/int64 uint8/uint16/uint32/uint64 float32/float64 string
// NOT support int/uint
type Set struct {
	m map[interface{}]bool
}

func New() *Set {
	return &Set{
		m: make(map[interface{}]bool),
	}
}

// if alrady existed, return false, otherwise return true
func (s *Set) Add(val interface{}) bool {
	_, exited := s.m[val]
	s.m[val] = true
	return !exited
}

func (s *Set) Delete(val interface{}) (interface{}, bool) {
	old, existed := s.m[val]
	delete(s.m, val)
	return old, existed
}

func (s *Set) Exists(val interface{}) bool {
	_, existed := s.m[val]
	return existed
}

func (s *Set) Inter(s2 *Set) *Set {
	s3 := New()
	for val := range s.m {
		s3.m[val] = true
	}
	for val := range s2.m {
		s3.m[val] = true
	}
	return s3
}

func (s *Set) Union(s2 *Set) *Set {
	s3 := New()
	for val := range s.m {
		if _, existed := s2.m[val]; existed {
			s3.m[val] = true
		}
	}
	return s3
}

func (s *Set) Diff(s2 *Set) *Set {
	s3 := New()
	for val := range s.m {
		if _, existed := s2.m[val]; !existed {
			s3.m[val] = true
		}
	}
	return s3
}

func (s *Set) Vals() []interface{} {
	var vals []interface{}
	for val := range s.m {
		vals = append(vals, val)
	}
	return vals
}

func (s *Set) Length() int {
	return len(s.m)
}

// Only for test
func (s *Set) String() string {
	var str string
	str += "{"

	var vals []string
	for val := range s.m {
		vals = append(vals, fmt.Sprintf("%v", val))
	}

	sort.Strings(vals)
	for _, val := range vals {
		str += fmt.Sprintf("%v, ", val)
	}

	str = str[:len(str)-2]
	str += "}"
	return str
}

func (s *Set) Marshal(w io.Writer) error {
	// write total. the total must > 0
	total := s.Length()
	if err := util.Write(w, int64(total)); err != nil {
		return err
	}

	for val := range s.m {
		// use Plain to marshal
		p := plain.New(val)
		if err := p.Marshal(w); err != nil {
			return err
		}
	}
	return nil
}

func (s *Set) Unmarshal(r io.Reader) error {
	var total int64
	if err := util.Read(r, &total); err != nil {
		return err
	}

	for i := 0; i < int(total); i++ {
		p := plain.New(nil)
		if err := p.Unmarshal(r); err != nil {
			return err
		}
		s.m[p.Val()] = true
	}
	return nil
}
