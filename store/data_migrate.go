package store

import (
	"sync"
	"uplus.io/udb/hash"
	"uplus.io/udb/proto"
)

const pageSize = 100

type DataMigrate struct {
	engine   *Engine
	dataComm proto.ClusterDataCommunication
	receiver int32

	queue []*proto.DataBody
	sync.Mutex
}

func NewDataMigrate(engine *Engine, dataComm proto.ClusterDataCommunication, receiver int32) *DataMigrate {
	return &DataMigrate{engine: engine, dataComm: dataComm, queue: make([]*proto.DataBody, 0), receiver: receiver}
}

//1 3 5 7 9 11
//2
func (p *DataMigrate) Migrate(startRing int32, endRing int32) {
	for i := 0; i < p.engine.PartSize(); i++ {
		store := p.engine.part(int32(i))
		//part := store.Part()
		store.MetaForEach(func(key []byte, meta proto.DataMeta) bool {
			p.Lock()
			defer p.Unlock()
			identity := NewIdentity(meta.Namespace, meta.Table, meta.Key)
			ring := hash.Int32(identity.IdBytes())
			if startRing < endRing {
				if ring >= startRing && ring < endRing {
					content, err := store.GetData(*identity)
					p.put(meta, *content, err)
				}
			} else {
				if ring >= startRing {
					content, err := store.GetData(*identity)
					p.put(meta, *content, err)
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

func (p *DataMigrate) put(meta proto.DataMeta, data proto.DataContent, err error) {
	dataBody := &proto.DataBody{}
	dataBody.Namespace = meta.Namespace
	dataBody.Table = meta.Table
	dataBody.Id = meta.Key
	dataBody.Version = meta.Version
	p.queue = append(p.queue, dataBody)

	if len(p.queue) >= pageSize {
		p.pushData()
	}
}

func (p *DataMigrate) pushData() {
	p.Lock()
	defer p.Unlock()
	response := p.dataComm.Push(p.receiver, &proto.PushRequest{Data: p.queue})
	if response != nil && response.Success {
		p.queue = make([]*proto.DataBody, 0)
	}
}
