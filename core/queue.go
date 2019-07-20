/*
 * Copyright (c) 2019 uplus.io
 */

package core

type LinkNode struct {
	data interface{}
	next *LinkNode
}

type Queue struct {
	head *LinkNode
	end  *LinkNode
}

func NewQueue() *Queue {
	q := &Queue{nil, nil}
	return q
}

func (q *Queue) Push(data interface{}) {
	n := &LinkNode{data: data, next: nil}

	if q.end == nil {
		q.head = n
		q.end = n
	} else {
		q.end.next = n
		q.end = n
	}

	return
}

func (q *Queue) Pop() (interface{}, bool) {
	if q.head == nil {
		return nil, false
	}

	data := q.head.data
	q.head = q.head.next
	if q.head == nil {
		q.end = nil
	}

	return data, true
}
