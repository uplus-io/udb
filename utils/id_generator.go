/*
 * Copyright (c) 2019 uplus.io
 */

package utils

import (
	"github.com/rs/xid"
	"github.com/satori/go.uuid"
)

//https://github.com/rs/xid
func GenId() string {
	return xid.New().String()
}

//https://github.com/satori/go.uuid
func GenUUID() string {
	return uuid.NewV4().String()
}
