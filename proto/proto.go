/*
 * Copyright (c) 2019 uplus.io
 */

package proto

import (
	ggproto "github.com/golang/protobuf/proto"
)

func Marshal(pb ggproto.Message) ([]byte, error) {
	return ggproto.Marshal(pb)
}

func Unmarshal(buf []byte, pb ggproto.Message) (error) {
	return ggproto.Unmarshal(buf, pb)
}
