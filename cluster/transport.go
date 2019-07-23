/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

import (
	"net"
)

type TransportStatus uint8

const (
	TransportStatusSuspect TransportStatus = iota
	TransportStatusAlive
	TransportStatusDead
	TransportStatusLeave
)

type TransportConfig struct {
	Id            int32
	BindIp        []string
	BindPort      int
	AdvertisePort int
	Seeds         []string
	Secret        string

	EventListener  EventListener
	PacketListener PacketListener
}

type TransportInfo struct {
	Id     int32          //节点Id
	Name   string          //节点名称
	Status TransportStatus //节点状态
	Addr   net.IP          //节点ip
	Port   uint16          //节点端口
	Native interface{}     //原生实现
}

type Transport interface {
	Serving() *TransportInfo
	Shutdown()
	SendToTCP(nodeId int32, msg []byte) error
	SendToUDP(nodeId int32, msg []byte) error
	Me() TransportInfo
}
