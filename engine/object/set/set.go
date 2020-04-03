package set

import (
	"fmt"
	"github.com/clovers4/gres/util"
	"io"
	"sort"
)

//todo

// only support int8/int16/int32/int64 uint8/uint16/uint32/uint64 float32/float64 string
// NOT support int/uint
type Set struct {
	m map[string]bool
}

func New() *Set {
	return &Set{
		m: make(map[string]bool),
	}
}

func (s *Set) Add(val string) {
	s.m[val] = true
}

func (s *Set) Delete(val string) {
	delete(s.m, val)
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
		vals = append(vals, val)
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
		if err := util.Write(w, val); err != nil {
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
		var val string
		if err := util.Read(r, &val); err != nil {
			return err
		}
		s.m[val] = true
	}
	return nil
}
