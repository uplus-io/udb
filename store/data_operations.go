package store

import (
	log "uplus.io/udb/logger"
	"sync"
	"uplus.io/udb"
	"uplus.io/udb/hash"
	"uplus.io/udb/proto"
)

const pageSize = 100

type DataOperations struct {
	engine   *Engine
	dataComm proto.ClusterDataCommunication
	receiver int32

	queue []*proto.DataBody
	sync.Mutex
}

func NewDataOperations(engine *Engine, dataComm proto.ClusterDataCommunication, receiver int32) *DataOperations {
	return &DataOperations{engine: engine, dataComm: dataComm, queue: make([]*proto.DataBody, 0), receiver: receiver}
}

//1 3 5 7 9 11
//2
func (p *DataOperations) Migrate(startRing int32, endRing int32) {
	for i := 0; i < p.engine.PartSize(); i++ {
		store := p.engine.part(int32(i))
		part := store.Part()
		store.MetaForEach(func(key []byte, meta proto.DataMeta) bool {
			p.Lock()
			defer p.Unlock()
			identity := NewIdentity(meta.Namespace, meta.Table, meta.Key)
			ring := hash.Int32(identity.IdBytes())
			if startRing < endRing {
				if ring >= startRing && ring < endRing {
					_, content, err := store.GetData(*identity)
					p.put(*part, meta, *content, err)
				}
			} else {
				if ring >= startRing {
					_, content, err := store.GetData(*identity)
					p.put(*part, meta, *content, err)
				}
			}
			return true
		})
	}
	if len(p.queue) > 0 {
		p.pushData()
	}
	p.dataComm.MigrateResponse(p.receiver, &proto.DataMigrateResponse{Completed: true})
}

func (p *DataOperations) put(partition proto.Partition, meta proto.DataMeta, data proto.DataContent, err error) {
	dataBody := &proto.DataBody{}
	dataBody.Namespace = meta.Namespace
	dataBody.Table = meta.Table
	dataBody.Id = meta.Key
	dataBody.PartitionId = partition.Id
	dataBody.PartitionIndex = partition.Index
	dataBody.Version = meta.Version
	dataBody.Ring = meta.Ring
	dataBody.Content = data.Content
	p.queue = append(p.queue, dataBody)

	if len(p.queue) >= pageSize {
		p.pushData()
	}
}

func (p *DataOperations) pushData() {
	p.Lock()
	defer p.Unlock()
	response := p.dataComm.Push(p.receiver, &proto.PushRequest{Data: p.queue})
	if response != nil && response.Success {
		p.queue = make([]*proto.DataBody, 0)
	}
}

// A节点将数据(key,value,version)及对应的版本号推送给B节点
// B节点更新A发过来的数据中比自己新的数据
func (p *DataOperations) Push(dataArray []*proto.DataBody) {
	result := make([]*proto.DataBody, len(dataArray))
	for i, data := range dataArray {
		identity := NewIdentity(data.GetNamespace(), data.GetTable(), data.GetId())
		part := p.MatchPartition(data.PartitionId, data.PartitionIndex, data.Ring)
		meta, _, err := p.engine.GetData(part.Id, *identity)
		if err == udb.ErrDbKeyNotFound || meta.Version < data.Version {
			dat := NewData(*NewIdentity(data.Namespace, data.Table, data.Id), data.Content)
			dat.Version = data.Version
			p.engine.SetData(part.Id, *dat, false)
			log.Debugf(
				"local data older,update[%s]ring:%d from v:%d to v:%d",
				string(identity.IdBytes()), uint32(data.Ring), meta.Version,data.Version)
		}
		result[i] = &proto.DataBody{Namespace: data.Namespace, Table: data.Table, Id: data.Id, Version: data.Version}
	}
	response := &proto.PushResponse{Success: true, Data: result}
	p.dataComm.PushReply(p.receiver, response)
}

// A不发送数据的value，仅发送数据的摘要key和version给B。
// B根据版本比较数据，将本地比A新的数据(key,value,version)推送给A
// A更新自己的本地数据
func (p *DataOperations) Pull([]*proto.DataBody) {

}

func (p *DataOperations) MatchPartition(partId, partIndex, ring int32) *proto.Partition {
	ranges := p.engine.PartRanges()
	length := len(ranges)
	if ranges[0] <= int(ring) && int(ring) <= ranges[length-1] {
		return p.engine.SimilarPart(ring)
	} else {
		return p.engine.GetPartOfIndex(partIndex)
	}
}
