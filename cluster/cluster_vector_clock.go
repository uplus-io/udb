/*
 * Copyright (c) 2019 uplus.io 
 */

package cluster

import (
	"sort"
	"sync"
)

type Clock struct {
	NodeId    uint32
	Remote    int64
	Local     int64
	Available bool
}

func NewClock(nodeId uint32, remote int64, local int64, available bool) *Clock {
	return &Clock{NodeId: nodeId, Remote: remote, Local: local, Available: available}
}

type Clocks struct {
	clocks []Clock

	sync.RWMutex
}

func (p *Clocks) Len() int {
	return len(p.clocks)
}

//时钟按照本地时间从低到高排序
func (p *Clocks) Less(i, j int) bool {
	return p.clocks[i].Local < p.clocks[j].Local
}
func (p *Clocks) Swap(i, j int) {
	tmp := p.clocks[i]
	p.clocks[i] = p.clocks[j]
	p.clocks[j] = tmp
}

func (p *Clocks) GetClock(nodeId uint32) (int, *Clock) {
	p.RLock()
	defer p.RUnlock()
	var lastIndex = 0
	for lastIndex, clock := range p.clocks {
		if clock.NodeId == nodeId {
			return lastIndex, &clock
		}
	}
	return lastIndex, nil
}

func (p *Clocks) UpdateClock(clock *Clock) {
	p.Lock()
	defer p.Unlock()
	index, find := p.GetClock(clock.NodeId)
	if find != nil {
		find = clock
	} else {
		p.clocks[index+1] = *clock
	}
	sort.Sort(p)
}

type Version struct {
	NodeId    uint32
	Version   uint32
	Timestamp int64
}

type Versions struct {
	Versions []Version
}

func (p *Versions) Len() int {
	return len(p.Versions)
}

//按版本号倒序排列
func (p *Versions) Less(i, j int) bool {
	return p.Versions[i].Version > p.Versions[j].Version
}
func (p *Versions) Swap(i, j int) {
	tmp := p.Versions[i]
	p.Versions[i] = p.Versions[j]
	p.Versions[j] = tmp
}

func (p *Versions) GetVersion(nodeId uint32) (int, *Version) {
	var lastIndex = 0
	for lastIndex, ver := range p.Versions {
		if ver.NodeId == nodeId {
			return lastIndex, &ver
		}
	}
	return lastIndex, nil
}

func (p *Versions) UpdateVersion(ver *Version) {
	index, find := p.GetVersion(ver.NodeId)
	if find != nil {
		find = ver
	} else {
		p.Versions[index+1] = *ver
	}
	sort.Sort(p)
}

type DataVectorClock struct {
	Clocks *Clocks
}

/**
前提：最低过半节点存活 集群有效 2/3
版本比较场景
每次更新至少更新过半节点才返回成功，否则返回失败

场景1：完全一致
a:1:100 b:1:100 c:1:100
a:1:100 b:1:100 c:1:100
场景2：过半一致
a:2:100 b:2:100 c:1:90
a:2:100 b:2:100 c:1:90
push c
场景3：版本半数冲突
a:2:101 b:2:100 c:1:90
a:2:100 b:2:102 c:1:90
取用下方数据 
 */

func (p *DataVectorClock) Compare(left Versions, right Versions) int {
	return 0
}
