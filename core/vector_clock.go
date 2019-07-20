/*
 * Copyright (c) 2019 uplus.io
 */

package core

type VectorClock struct {
	//机器标志
	Machine uint32 `json:"m"`
	//该机器上的版本
	Version uint32 `json:"v"`
	//该机器上的时间戳
	Timestamp int64 `json:"t"`
}

type VectorVersion struct {
	//机器总数
	Size int
	//各机器相较于本机的时间戳
	Timestamps []uint64
	//机器向量时钟
	VectorClocks []VectorClock
}

func NewVectorVersion() *VectorVersion {
	version := &VectorVersion{}
	return version
}

func (p *VectorVersion) Update(machine uint32, vc VectorClock) {

}
