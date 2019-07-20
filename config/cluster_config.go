/*
 * Copyright (c) 2019 uplus.io
 */

package config

type ClusterConfig struct {
	Name           string         `json:"level" yaml:"name"`
	Seeds          []string       `json:"seeds" yaml:"seeds"`
	BindIp         []string       `json:"bind_ip" yaml:"bind_ip"`
	BindPort       int            `json:"bind_port" yaml:"bind_port"`
	AdvertisePort  int            `json:"advertise_port" yaml:"advertise_port"`
	ReplicaCount   int            `json:"replica_count" yaml:"replica_count"`
	SecurityConfig SecurityConfig `json:"security" yaml:"security"`
	StorageConfig  StorageConfig  `json:"storage" yaml:"storage"`
	LogConfig      LogConfig      `json:"log" yaml:"log"`
}

type SecurityConfig struct {
	Secret string `json:"secret" yaml:"secret"`
}

type StorageConfig struct {
	Engine     string   `json:"engine" yaml:"engine"`
	Meta       string   `json:"meta" yaml:"meta"`
	Wal        string   `json:"wal" yaml:"wal"`
	Partitions []string `json:"partitions" yaml:"partitions"`
}

type LogConfig struct {
	Path     string `json:"path" yaml:"path"`
	Filename string `json:"filename" yaml:"filename"`
	Level    uint32 `json:"level" yaml:"level"`
}
