package store

import "uplus.io/udb/proto"

type DataOperations interface {
	Push(dataArray []*proto.DataBody) error
	Pull([]*proto.DataBody) error
}

type DataPipeline interface {

}
