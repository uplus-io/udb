/*
 * Copyright (c) 2019 uplus.io
 */

package deprecated

type TableType uint8

const (
	TableTypeKV TableType = iota

)

type Descriptor struct {
	//空间
	Namespace string `json:"namespace"`
	//表名
	Table string `json:"table"`
}
