package list

import (
	"fmt"
	"github.com/clovers4/gres/engine/object/plain"
	"github.com/clovers4/gres/util"
	"io"
)

type node struct {
	prev *node
	next *node

	val interface{}
}

func newNode(val interface{}) *node {
	return &node{val: val}
}

func (n *node) Next() *node {
	return n.next
}
func (n *node) Prev() *node {
	return n.prev
}

func (n *node) SetVal(val interface{}) interface{} {
	old := n.val
	n.val = val
	return old
}

func (n *node) Val() interface{} {
	return n.val
}

type List struct {
	header *node
	tail   *node

	length int
}

func New() *List {
	return new(List)
}

func (ls *List) LPush(val interface{}) {
	n := newNode(val)
	n.next = ls.header
	if ls.length == 0 {
		ls.tail = n
	} else {
		ls.header.prev = n
	}
	ls.header = n
	ls.length++
}

// NOTICE: LPop assume the length > 0, so cannot distinguish nil or nil node
func (ls *List) LPop() interface{} {
	if ls.length == 0 {
		return nil
	}

	old := ls.header
	ls.header = ls.header.next
	if ls.length == 1 {
		ls.tail = nil
	} else {
		ls.header.prev = nil
	}
	ls.length--
	return old.val
}

func (ls *List) RPush(val interface{}) {
	n := newNode(val)
	n.prev = ls.tail
	if ls.length == 0 {
		ls.header = n
	} else {
		ls.tail.next = n
	}
	ls.tail = n
	ls.length++
}

// NOTICE: RPop assume the length > 0, so cannot distinguish nil or nil node
func (ls *List) RPop() interface{} {
	if ls.length == 0 {
		return nil
	}

	n := ls.tail
	ls.tail = ls.tail.prev
	if ls.length == 1 {
		ls.header = nil
	} else {
		ls.tail.next = nil
	}
	ls.length--
	return n.val
}

func (ls *List) Front() *node {
	return ls.header
}

func (ls *List) End() *node {
	return ls.tail
}

// index start at 0
func (ls *List) Index(index int) *node {
	if index < 0 {
		index = ls.Length() + index
	}

	if index < 0 || index >= ls.Length() {
		return nil
	}

	// 从左向右遍历, 否则从右向左
	var n *node
	if index <= ls.Length()/2 {
		n = ls.header
		for i := 0; i < index; i++ {
			n = n.Next()
		}
	} else {
		n = ls.tail
		for i := ls.Length() - 1; i > index; i-- {
			n = n.Prev()
		}
	}
	return n
}

func (ls *List) Length() int {
	return ls.length
}

// Only for test
func (ls *List) String() string {
	var s string
	s += "{"
	for n := ls.Front(); n != nil; n = n.Next() {
		s += fmt.Sprintf("%v, ", n.Val())
	}
	if len(s) > 2 {
		s = s[:len(s)-2]
	}
	s += "}"
	return s
}

func (ls *List) Marshal(w io.Writer) error {
	// write total. the total must > 0
	total := ls.Length()
	if err := util.Write(w, int64(total)); err != nil {
		return err
	}

	// loop write score and val
	for n := ls.Front(); n != nil; n = n.Next() {
		// use Plain to marshal
		p := plain.New(n.Val())
		if err := p.Marshal(w); err != nil {
			return err
		}
	}
	return nil
}

func (ls *List) Unmarshal(r io.Reader) error {
	var total int64
	if err := util.Read(r, &total); err != nil {
		return err
	}

	for i := 0; i < int(total); i++ {
		p := plain.New(nil)
		if err := p.Unmarshal(r); err != nil {
			return err
		}

		ls.RPush(p.Val())
	}
	return nil
}
