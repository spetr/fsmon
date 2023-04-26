package main

import (
	"os"
	"time"

	"github.com/kardianos/service"
	"github.com/spetr/go-zabbix-sender"
)

var logger service.Logger

type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Started in terminal.")
	} else {
		logger.Info("Started with service manager.")
	}
	p.exit = make(chan struct{})
	go p.run()
	return nil
}

func (p *program) run() {
	if err := configLoad(); err != nil {
		logger.Error("Config load error: ", err)
		os.Exit(1)
	}

	// Run zabbix discovery
	metrics := make([]*zabbix.Metric, 1)
	t := time.Now().Unix()
	metrics[0] = zabbix.NewMetric(conf.Zabbix.Hostname, "fsmon.discovery", "[{\"{#NAME}\": \"name\"}]", true, t)
	for i := range conf.Zabbix.Servers {
		logger.Infof("Starting service discovery for host %s on %s", conf.Zabbix.Hostname, conf.Zabbix.Servers[i].Host)
		zabbixSender := zabbix.NewSender(conf.Zabbix.Servers[i].Host)
		zabbixSender.ConnectTimeout = conf.Zabbix.Servers[i].ConnectTimeout
		zabbixSender.ReadTimeout = conf.Zabbix.Servers[i].ReadTimeout
		zabbixSender.WriteTimeout = conf.Zabbix.Servers[i].WriteTimeout
		zabbixResponse, err, _, _ := zabbixSender.SendMetrics(metrics)
		if *debugFlag {
			logger.Infof("Zabbix response info: %s", zabbixResponse.Info)
			logger.Infof("Zabbix response: %s", zabbixResponse.Response)
		}
		if err != nil {
			logger.Errorf("Failed to send metrics to %s: %s", conf.Zabbix.Servers[i].Host, err)
		}
	}

	// Start Zabbix sender
	for i := range conf.Filesystems {
		logger.Infof("Starting monitoring of %s (%s)", conf.Filesystems[i].Name, conf.Filesystems[i].Mountpoint)
		go monFsUpdate(conf.Filesystems[i].Mountpoint, conf.Filesystems[i].Name)
	}

	<-p.exit
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stopping")
	close(p.exit)
	return nil
}
