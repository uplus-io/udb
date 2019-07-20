/*
 * Copyright (c) 2019 uplus.io
 */

package core

import (
	"github.com/hashicorp/memberlist"
	"net"
)

type NodeStatus uint8

const (
	NodeStatusSuspect NodeStatus = iota
	NodeStatusAlive
	NodeStatusDead
	NodeStatusLeave
)

type Node struct {
	//节点id
	Id     uint32
	Status NodeStatus
	//节点ip
	Addr net.IP
	//节点端口
	Port      uint16
	//Timestamp proto.Timestamp
	Native    *memberlist.Node
}
