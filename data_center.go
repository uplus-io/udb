/*
 * Copyright (c) 2019 uplus.io
 */

package udb

type DataCenter struct {
	Id    uint32
	Size  uint32
	Racks map[uint32]*Rack
	Next  *DataCenter
}
