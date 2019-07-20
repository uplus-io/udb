/*
 * Copyright (c) 2019 uplus.io
 */

package core

import "sort"

type ArrayItem interface {
	GetId() int32
	Compare(item ArrayItem) int32
}

type Array struct {
	Collection []ArrayItem
	tail       int
	length     int
}

func NewArray() *Array {
	return &Array{Collection: make([]ArrayItem, 0, 3)}
}

func (p *Array) Len() int {
	return p.length
}

func (p *Array) Less(i, j int) bool {
	return p.Collection[i].Compare(p.Collection[j]) < 0
}

func (p *Array) Swap(i, j int) {
	tmp := p.Collection[i]
	p.Collection[i] = p.Collection[j]
	p.Collection[j] = tmp
}

func (p *Array) Add(item ArrayItem) ArrayItem {
	exist := p.Unique(item)
	if exist == nil {
		exist = item
		if p.tail >= len(p.Collection) {
			p.Collection = append(p.Collection, item)
		} else {
			p.Collection[p.tail] = item
		}
		p.tail++
		p.length++
		p.Collection = append(p.Collection, item)

		sort.Sort(p)
	} else {
		item = exist
	}
	return exist
}

func (p *Array) Delete(id int32) bool {
	for i := 0; i < p.length; i++ {
		item := p.Collection[i]
		if item.GetId() == id {
			p.Collection = append(p.Collection[:i], p.Collection[i+1:]...)
			sort.Sort(p)
			return true
		}
	}
	return false
}

func (p *Array) IfAbsentCreate(item ArrayItem) ArrayItem {
	exist := p.Unique(item)
	if exist == nil {
		return p.Add(item)
	}
	return exist
}

func (p *Array) Exists(item ArrayItem) bool {
	return p.Unique(item) != nil
}

func (p *Array) NotExists(item ArrayItem) bool {
	return p.Unique(item) == nil
}

func (p *Array) Index(index int) ArrayItem {
	return p.Collection[index]
}

func (p *Array) Id(id int32) ArrayItem {
	for i := 0; i < p.length; i++ {
		item := p.Collection[i]
		if item.GetId() == id {
			return item
		}
	}
	return nil
}

func (p *Array) Unique(input ArrayItem) ArrayItem {
	for i := 0; i < p.length; i++ {
		item := p.Collection[i]
		if item.Compare(input) == 0 {
			return item
		}
	}
	return nil
}
