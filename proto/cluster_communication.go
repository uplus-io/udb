package proto

type ClusterCommunication interface {
	SendNodeInfoTo(to int32) error
}

type ClusterDataCommunication interface {
	Push(to int32, request *PushRequest) *PushResponse
	Pull(to int32, request *PullRequest) *PullRequest
	MigrateRequest(to int32, request *DataMigrateRequest) error
	MigrateResponse(to int32, request *DataMigrateResponse) error
}
