package proto

func NewDescription(ns, tab int32) *Description {
	return &Description{Namespace: ns, Table: tab}
}

