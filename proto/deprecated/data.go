/*
 * Copyright (c) 2019 uplus.io
 */

package deprecated

type Data struct {
	//字段名
	Descriptor string `json:"descriptor"`
	//版本
	Version uint32 `json:"version"`
	//哈希值
	Hashcode uint32 `json:"hashcode"`
	//时间戳
	Timestamp int64 `json:"timestamp"`
	//键名
	Key []byte
	//值
	Value []byte
}

type KVData struct {
	//键名
	Key []byte
	//值
	Value []byte
}
