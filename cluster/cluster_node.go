/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

import (
	"uplus.io/udb/proto"
)

type NodeEventType uint8

const (
	NodeEventType_Join NodeEventType = iota
	NodeEventType_Leave
	NodeEventType_Update
)

type NodeEvent struct {
	Type   NodeEventType
	Node   *proto.Node
	Native interface{}
}

func NewNodeEvent(t NodeEventType, node *proto.Node, native interface{}) *NodeEvent {
	return &NodeEvent{Type: t, Node: node, Native: native}
}
