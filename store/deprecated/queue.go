/*
 * Copyright (c) 2019 uplus.io
 */

package deprecated

type QueueData interface {
}

type QueueNode struct {
	Item *QueueData
	Next *Queue
}

type Queue struct {
	Size  uint32
	First *QueueNode
	Last  *QueueNode
}
