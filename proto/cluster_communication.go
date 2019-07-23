package proto

type ClusterCommunication interface {
	SendNodeInfoTo(to int32) error
}
