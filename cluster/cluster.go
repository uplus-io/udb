/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

import (
	"sync"
	"uplus.io/udb/config"
	"uplus.io/udb/data"
	"uplus.io/udb/logger"
	log "uplus.io/udb/logger"
	"uplus.io/udb/store"
)

type Cluster struct {
	config config.ClusterConfig

	engine    *store.Engine
	warehouse *data.Warehouse
	exit      chan bool

	transport Transport
	pipeline  Pipeline

	node *data.Node

	launched bool

	sync.RWMutex
}

func NewCluster(config config.ClusterConfig) *Cluster {
	//init log
	logger.DebugLoggerEnable(true)
	logConfig := config.LogConfig
	level := logger.LoggerLevel(logConfig.Level)
	logger.NewLogger(level, logConfig.Path, logConfig.Filename)
	return &Cluster{
		config:    config,
		warehouse: data.NewWarehouse(),
		pipeline:  NewPipelinePacket(),
		exit:      make(chan bool),
	}
}

func (p *Cluster) Listen() {
	p.startEngine()
	p.launchGossip()
	p.packetLoop()
}

func (p *Cluster) startEngine() {
	p.engine = store.NewEngine(p.config.StorageConfig)
}

func (p *Cluster) launchGossip() {
	p.Lock()
	defer p.Unlock()
	p.transport = NewTransportGossip(
		&TransportConfig{
			Seeds:          p.config.Seeds,
			Secret:         p.config.SecurityConfig.Secret,
			BindIp:         p.config.BindIp,
			BindPort:       p.config.BindPort,
			AdvertisePort:  p.config.AdvertisePort,
			EventListener:  NewClusterEventListener(p.warehouse),
			PacketListener: NewClusterPacketListener(p.pipeline)})
	transportInfo := p.transport.Serving()
	p.launched = true
	log.Debugf("cluster node[%d] started %v", transportInfo.Id, p.launched)
	//p.initNodes()
}
func (p *Cluster) packetLoop() {
	for {
		packet := <-p.pipeline.In()
		node := p.warehouse.GetNode(uint32(packet.GetDataCenter()), uint32(packet.GetTo()))
		if node != nil {
		}
	}
}
