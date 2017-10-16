package graph

import "fmt"

type Id interface {
	Value() interface{}
	String() string
}

func NewId(i interface{}) Id {
	return &id{i}
}

type id struct {
	val interface{}
}

func (i *id) Value() interface{} {
	return i.val
}

func (i *id) String() string {
	return fmt.Sprintf("Id(%v)", i.val)
}
