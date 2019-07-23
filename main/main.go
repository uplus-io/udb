/*
 * Copyright (c) 2019 uplus.io
 */

package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"uplus.io/udb/cluster"
	"uplus.io/udb/config"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		panic("must give config file path")
	}

	configPath := args[1]

	fmt.Printf("uplus db ready read config file[%s]\n", configPath)

	file, e := os.OpenFile(configPath, os.O_RDONLY, 0600)
	if e != nil {
		panic("open config file fail")
	}
	config := config.ClusterConfig{}
	yaml.NewDecoder(file).Decode(&config)

	cluster := cluster.NewCluster(config)
	cluster.Listen()
}
