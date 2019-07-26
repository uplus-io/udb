package proto

import (
	"strings"
	"uplus.io/udb/utils"
)

func NewDescription(ns, tab int32) *Description {
	return &Description{Namespace: ns, Table: tab}
}

func ParseRepository(ip string) *Repository {
	bits := strings.Split(ip, ".")
	center := utils.StringToInt32(bits[0])
	area := utils.StringToInt32(bits[1])
	rack := utils.StringToInt32(bits[2])
	return &Repository{DataCenter: center, Area: area, Rack: rack}
}
