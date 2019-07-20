/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

type EndpointSnitch interface {
	DataCenter() uint32

	Rack() uint32
}
