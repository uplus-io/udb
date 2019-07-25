/*
 * Copyright (c) 2019 uplus.io
 */

package core

import (
	"net"
)

type Node struct {
	//节点id
	Id int32
	//节点ip
	Addr net.IP
	//节点端口
	Port int32
	//Timestamp proto.Timestamp
	Native interface{}
}
