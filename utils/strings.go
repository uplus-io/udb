/*
 * Copyright (c) 2019 uplus.io
 */

package utils

import "bytes"

func StringEqual(s1 string, s2 string) bool {
	return bytes.Equal([]byte(s1), []byte(s2))
}
