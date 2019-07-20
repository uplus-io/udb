/*
 * Copyright (c) 2019 uplus.io
 */

package proto

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"uplus.io/udb/core"
)

func NewNode(id int32, ip string, port int32) *Node {
	return &Node{Id: id, Ip: ip, Port: port}
}

func (p *Node) Compare(item core.ArrayItem) int32 {
	return p.GetId() - item.GetId()
}

func (v Node) Address() string {
	return fmt.Sprintf("%s:%d", v.Ip, v.Port)
}

func (v Node) Addr() net.IP {
	bits := strings.Split(v.Ip, ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	return net.IPv4(byte(b0), byte(b1), byte(b2), byte(b3))
}
