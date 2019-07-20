/*
 * Copyright (c) 2019 uplus.io
 */

package core

type Serializable interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) (error)
}
