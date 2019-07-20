/*
 * Copyright (c) 2019 uplus.io
 */

package udb

type BaseStream struct {
	path string

	//streams []Stream
}

func NewBaseStream(path string) *BaseStream {
	baseStream := &BaseStream{path: path}
	return baseStream
}

//func (p *BaseStream) init() {
	//streams := make([]Stream, 4)
	//streams[0] = newStream("2.test.dat", 4+8+2)
	//streams[1] = newStream("8.test.dat", 4+8+8)
	//streams[2] = newStream("1k.test.dat", 4+8+1024)
	//streams[3] = newStream("32k.test.dat", 4+8+3*1024)
	//p.streams = streams
//}

//func (p *BaseStream) WriteNumber(val interface{}) (*Base, error) {
//	var base *Base
//
//	return base, nil
//}
