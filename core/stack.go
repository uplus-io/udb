/*
 * Copyright (c) 2019 uplus.io
 */

package core

type Stack struct {
	head *LinkNode
}

func NewStack() *Stack {
	s := &Stack{nil}
	return s
}

func (s *Stack) Push(data interface{}) {
	n := &LinkNode{data: data, next: s.head}
	s.head = n
}

func (s *Stack) Pop() (interface{}, bool) {
	n := s.head
	if s.head == nil {
		return nil, false
	}
	s.head = s.head.next
	return n.data, true
}
