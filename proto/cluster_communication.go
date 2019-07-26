package proto

type ClusterCommunication interface {
	SendNodeInfoTo(to int32) error
}

type ClusterDataCommunication interface {
	Push(to int32, request *PushRequest) *PushResponse
	PushReply(to int32, request *PushResponse) error
	Pull(to int32, request *PullRequest) *PullResponse
	PullReply(to int32, request *PullResponse) error
	MigrateRequest(to int32, request *DataMigrateRequest) error
	MigrateResponse(to int32, request *DataMigrateResponse) error
}
