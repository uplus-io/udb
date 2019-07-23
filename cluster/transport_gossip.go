/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

import (
	"fmt"
	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/memberlist"
	"uplus.io/udb"
	"uplus.io/udb/hash"
	log "uplus.io/udb/logger"
	"uplus.io/udb/proto"
	"uplus.io/udb/utils"
)

type TransportGossip struct {
	config  *TransportConfig
	members *memberlist.Memberlist
	nodes   map[int32]*memberlist.Node
	me      *TransportInfo
}

func NewTransportGossip(config *TransportConfig) *TransportGossip {
	return &TransportGossip{config: config, nodes: make(map[int32]*memberlist.Node)}
}

func CreateProtoNode(node *memberlist.Node) *proto.Node {
	id := utils.StringToInt32(node.Name)
	return proto.NewNode(id, node.Addr.String(), int32(node.Port))
}

func (p *TransportGossip) Serving() *TransportInfo {
	cfg := p.config
	if cfg.Id == 0 {
		var host string
		if cfg.BindIp[0] == "0.0.0.0" {
			ip, _ := sockaddr.GetPrivateIP()
			host = fmt.Sprintf("%s:%d", ip, cfg.BindPort)
		} else {
			host = fmt.Sprintf("%s:%d", cfg.BindIp, cfg.BindPort)
		}
		id := hash.Int32Of(host)
		cfg.Id = id
	}

	config := memberlist.DefaultLocalConfig()
	config.Name = utils.Int32ToString(cfg.Id)
	config.SecretKey = []byte(cfg.Secret)
	config.BindPort = cfg.BindPort
	config.AdvertisePort = cfg.AdvertisePort
	config.EnableCompression = true
	config.Delegate = p
	config.Events = p

	members, err := memberlist.Create(config)
	if err != nil {
		panic(fmt.Sprintf("Cluster lanuch fail:%v", err.Error()))
	}
	members.Join(cfg.Seeds)
	p.members = members
	p.me = &TransportInfo{Id: int32(cfg.Id),}

	for _, node := range members.Members() {
		p.NotifyJoin(node)
	}
	return p.me
}

func (p *TransportGossip) Shutdown() {
	p.members.Shutdown()
}

func (p *TransportGossip) Me() TransportInfo {
	return *p.me
}

func (p *TransportGossip) nativeNode(nodeId int32) *memberlist.Node {
	return p.nodes[nodeId]
}

func (p *TransportGossip) SendToTCP(nodeId int32, msg []byte) error {
	node := p.nativeNode(nodeId)
	if node != nil {
		return p.members.SendReliable(node, msg)
	}
	return udb.ErrClusterNodeOffline
}

func (p *TransportGossip) SendToUDP(nodeId int32, msg []byte) error {
	node := p.nativeNode(nodeId)
	if node != nil {
		return p.members.SendBestEffort(node, msg)
	}
	return udb.ErrClusterNodeOffline
}

//event

// NotifyJoin is invoked when a node is detected to have joined.
// The Node argument must not be modified.
func (p *TransportGossip) NotifyJoin(n *memberlist.Node) {
	node := CreateProtoNode(n)
	p.nodes[node.Id] = n
	p.config.EventListener.OnTopologyChanged(
		NewNodeEvent(NodeEventType_Join, node, n))
}

// NotifyLeave is invoked when a node is detected to have left.
// The Node argument must not be modified.
func (p *TransportGossip) NotifyLeave(n *memberlist.Node) {
	node := CreateProtoNode(n)
	delete(p.nodes, node.Id)
	p.config.EventListener.OnTopologyChanged(
		NewNodeEvent(NodeEventType_Leave, node, n))
}

// NotifyUpdate is invoked when a node is detected to have
// updated, usually involving the meta data. The Node argument
// must not be modified.
func (p *TransportGossip) NotifyUpdate(n *memberlist.Node) {
	node := CreateProtoNode(n)
	p.config.EventListener.OnTopologyChanged(
		NewNodeEvent(NodeEventType_Update, node, n))
}

//packet

// NodeMeta is used to retrieve meta-data about the current node
// when broadcasting an alive message. It's length is limited to
// the given byte size. This metadata is available in the Node structure.
func (p *TransportGossip) NodeMeta(limit int) []byte {
	//log.Debugf("NodeMeta [%d]\n", limit)
	return nil
}

// NotifyMsg is called when a user-data message is received.
// Care should be taken that this method does not block, since doing
// so would block the entire UDP packet receive loop. Additionally, the byte
// slice may be modified after the call returns, so it should be copied if needed
func (p *TransportGossip) NotifyMsg(dat []byte) {
	packet := &proto.Packet{}
	proto.Unmarshal(dat, packet)
	p.config.PacketListener.OnReceive(packet)
}

// GetBroadcasts is called when user data messages can be broadcast.
// It can return a list of buffers to send. Each buffer should assume an
// overhead as provided with a limit on the total byte size allowed.
// The total byte size of the resulting data to send must not exceed
// the limit. Care should be taken that this method does not block,
// since doing so would block the entire UDP packet receive loop.
func (p *TransportGossip) GetBroadcasts(overhead, limit int) [][]byte {
	//log.Debugf("GetBroadcasts [%d/%d]\n", overhead, limit)
	return nil
}

// LocalState is used for a TCP push/pull. This is sent to
// the remote side in addition to the membership information. Any
// data can be sent here. See MergeRemoteState as well. The `join`
// boolean indicates this is for a join instead of a push/pull.
func (p *TransportGossip) LocalState(join bool) []byte {
	//log.Debugf("LocalState [%v]\n", join)
	return nil
}

// MergeRemoteState is invoked after a TCP push/pull. This is the
// state received from the remote side and is the result of the
// remote side's LocalState call. The 'join'
// boolean indicates this is for a join instead of a push/pull.
func (p *TransportGossip) MergeRemoteState(buf []byte, join bool) {
	log.Debugf("MergeRemoteState [%s|%v]\n", string(buf), join)
}
