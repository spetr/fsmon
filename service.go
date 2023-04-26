package main

import (
	"os"

	"github.com/kardianos/service"
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
