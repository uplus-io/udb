/*
 * Copyright (c) 2019 uplus.io
 */

package deprecated

type BaseType uint8

const (
	UNKNOW_TYPE BaseType = iota
	//short data
	BASE_TYPE_BYTE
	BASE_TYPE_CHAR
	BASE_TYPE_BOOL
	//number data
	BASE_TYPE_SHORT
	BASE_TYPE_INT
	BASE_TYPE_LONG
	BASE_TYPE_FLOAT
	BASE_TYPE_DOUBLE
	//date data
	BASE_TYPE_DATE
	BASE_TYPE_TIME
	BASE_TYPE_DATETIME
	//STRING
	BASE_TYPE_STRING
	BASE_TYPE_TEXT
	BASE_TYPE_JSON
	//file
	BASE_TYPE_FILE
	//OBJ
	BASE_TYPE_OBJECT
)

type Base struct {
	BaseType BaseType
	length   int64
	//Buffer   *Buffer

	Page  int
	Index int
}

func NewBase(dataType BaseType) *Base {
	return &Base{BaseType: dataType}
}

func (b Base) HeadLength() int {
	return 4 + 8
}

func (b Base) Length() int {
	return b.HeadLength() + b.StorageLength()
}

func (b Base) StorageLength() int {
	baseType := b.BaseType
	if baseType >= BASE_TYPE_BYTE && baseType <= BASE_TYPE_BOOL {
		return 2
	} else if baseType >= BASE_TYPE_SHORT && baseType <= BASE_TYPE_DOUBLE {
		return 8
	} else if baseType >= BASE_TYPE_DATE && baseType <= BASE_TYPE_DATETIME {
		return 8
	} else if baseType >= BASE_TYPE_STRING && baseType <= BASE_TYPE_JSON {
		return 1024
	} else if baseType == BASE_TYPE_FILE {
		return 32 * 1024
	}
	return 0
}

//func (p Base) writeHead() {
//	p.Buffer.WriteInt8(int8(p.BaseType))
//	p.Buffer.WriteInt(p.Length())
//}
//
//func (p Base) WriteNumber(val interface{}) error {
//	buffer := NewBuffer(p.Length())
//	p.Buffer = buffer
//	p.writeHead()
//	switch val.(type) {
//	case int:
//		buffer.WriteInt(val.(int))
//	case int32:
//		buffer.WriteInt32(val.(int32))
//	case int64:
//		buffer.WriteLong(val.(int64))
//	case float32:
//	case float64:
//	default:
//		return errors.New("不支持的数据类型")
//	}
//	return nil
//}
