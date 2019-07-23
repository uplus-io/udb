/*
 * Copyright (c) 2019 uplus.io
 */

package udb

import "errors"

//cluster

var (
	ErrClusterNodeOffline = errors.New("cluster node offline")
)

var (
	ErrNotFoundClusterNode = errors.New("not found cluster node")

	ErrMessageDispatcherExist = errors.New("message dispatcher already exist")
	ErrMessageHandlerExist    = errors.New("message handler already exist")
)

var (
	ErrSerialize   = errors.New("serialize fail")
	ErrDeserialize = errors.New("deserialize fail")
)

//database
var (
	ErrPartNotFound    = errors.New("partition not found")
	ErrPartAllocated   = errors.New("partition has been allocated")
	ErrPartNotAllocate = errors.New("partition not allocate")
)

var (
	ErrDbKeyNotFound = errors.New("db:key not found")
)

//warehouse errors
