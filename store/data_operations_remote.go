package store

import "uplus.io/udb/proto"

type RemoteDataOperations struct {
}

func (p *RemoteDataOperations) Push(dataArray []*proto.DataBody) error {
	return nil
}
func (p *RemoteDataOperations) Pull([]*proto.DataBody) error {
	return nil
}
