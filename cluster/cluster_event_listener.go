/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

import (
	"uplus.io/udb/data"
	log "uplus.io/udb/logger"
)

type EventListener interface {
	OnTopologyChanged(event *NodeEvent)
}

type ClusterEventListener struct {
	warehouse *data.Warehouse
}

func NewClusterEventListener(warehouse *data.Warehouse) *ClusterEventListener {
	return &ClusterEventListener{warehouse: warehouse}
}

func (p *ClusterEventListener) OnTopologyChanged(event *NodeEvent) {
	switch event.Type {
	case NodeEventType_Join:
		p.warehouse.JoinNode(event.Node.Ip, int(event.Node.Port))
	case NodeEventType_Leave:
		p.warehouse.LeaveNode(event.Node.Ip, int(event.Node.Port))
	case NodeEventType_Update:
		log.Debugf("node[%s] event update", event.Node.String())
	}
}
