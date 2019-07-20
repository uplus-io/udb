/*
 * Copyright (c) 2019 uplus.io
 */

package deprecated

type Table struct {
	Id   uint32
	Name string
}

type TableRow struct {
	Id      uint32
	Key     interface{}
	Columns []TableColumn
}

type TableColumn struct {
	Name  string
	Value interface{}
}
