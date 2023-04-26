package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type (
	tConf struct {
		Filesystems []struct {
			Mountpoint string `yaml:"mountpoint"`
			Name       string `yaml:"name"`
		} `yaml:"filesystems"`
		Zabbix struct {
			Hostname string `yaml:"hostname"`
			Servers  []struct {
				Host           string        `yaml:"host"`
				FallbackDir    int           `yaml:"port"`
				ConnectTimeout time.Duration `yaml:"connectTimeout"`
				ReadTimeout    time.Duration `yaml:"readTimeout"`
				WriteTimeout   time.Duration `yaml:"writeTimeout"`
			} `yaml:"servers"`
		}
	}
)

var (
	conf tConf
)

func configLoad() (err error) {
	var (
		file *os.File
	)

	if file, err = os.Open(*configFile); err != nil {
		return
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(&conf); err != nil {
		return
	}

	if len(conf.Filesystems) == 0 {
		err = fmt.Errorf("No filesystems defined in config file")
		return
	}

	if conf.Zabbix.Hostname == "" {
		conf.Zabbix.Hostname, _ = os.Hostname()
	}

	return nil
}
