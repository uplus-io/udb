/*
 * Copyright (c) 2019 uplus.io
 */

package utils

import "strconv"

func StringToInt(str string) (int) {
	i, _ := strconv.ParseInt(str, 10, 32)
	return int(i)
}

func StringToInt32(str string) (int32) {
	i, _ := strconv.ParseInt(str, 10, 32)
	return int32(i)
}

func StringToInt64(str string) (int64) {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func IntToString(i int) (string) {
	return strconv.FormatInt(int64(i), 10)
}

func Int32ToString(i int32) (string) {
	return strconv.FormatInt(int64(i), 10)
}

func Int64ToString(i int64) (string) {
	return strconv.FormatInt(i, 10)
}
