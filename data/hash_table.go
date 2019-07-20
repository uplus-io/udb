/*
 * Copyright (c) 2019 uplus.io
 */

package data

import "uplus.io/udb/hash"

type HashTableData interface {
	Key() string
}

type HashTableNode struct {
	//256
	Key string
	//4
	Hash uint32
}

func NewHashTableNode(key string, hash uint32) *HashTableNode {
	return &HashTableNode{Key: key, Hash: hash}
}

type HashTable struct {
	n  int
	m  int
	st []HashTableNode
}

func NewHashTable(initialSize int) *HashTable {
	table := &HashTable{}
	table.m = initialSize
	table.st = make([]HashTableNode, initialSize)
	return table
}

//private void resize(int chains) {
//SeparateChainingHashST<Key, Value> temp = new SeparateChainingHashST<Key, Value>(chains);
//for (int i = 0; i < m; i++) {
//for (Key key : st[i].keys()) {
//temp.put(key, st[i].get(key));
//}
//}
//this.m  = temp.m;
//this.n  = temp.n;
//this.st = temp.st;
//}

func (p *HashTable) Size() int {
	return p.n
}

func (p *HashTable) IsEmpty() bool {
	return p.Size() == 0
}

func (p *HashTable) Contains(data HashTableData) bool {
	return false
}

func (p *HashTable) Get(data HashTableData) uint32 {
	hash := p.hash(data)
	node := p.st[hash]
	return node.Hash
}

func (p *HashTable) hash(data HashTableData) uint32 {
	hash := hash.UInt32Of(data.Key())
	return hash & 0x7fffffff % uint32(p.m)
}

/**
public void put(Key key, Value val) {
        if (key == null) throw new IllegalArgumentException("first argument to put() is null");
        if (val == null) {
            delete(key);
            return;
        }

        // double table size if average length of list >= 10
        if (n >= 10*m) resize(2*m);

        int i = hash(key);
        if (!st[i].contains(key)) n++;
        st[i].put(key, val);
    }
 */

func (p *HashTable) Put(data HashTableData) uint32 {
	//hash := p.hash(data)
	//if(p)
	//key := data.Key()
	//switch key.(type) {
	//case string:
	//	k := key.(string)
	//	hash := hash.UInt32Of(k)
	//	node := NewHashTableNode(k, hash)
	//	p.Table[p.Tail] = *node
	//	p.Tail++
	//}
	return 0
}

func Get(key string) {

}
