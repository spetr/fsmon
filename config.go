package main

import (
	"os"
	"strings"
)

type (
	tConf struct {
		Filesystem   string
		ZabbixServer []string
		Prometheus   string
	}
)

var conf tConf

func configLoad() (err error) {
	conf.Filesystem = os.Getenv("FS")
	conf.ZabbixServer = strings.Split(os.Getenv("ZABBIX_SERVER"), ",")
	conf.Prometheus = os.Getenv("PROMETHEUS")

	if conf.Filesystem == "" {
		conf.Filesystem = "/"
	}
	return nil
}
