package main

import (
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
	configLoad()

	for _, fs := range conf.Filesystems {
		go monFsUpdate(fs.Mountpoint, fs.Name)
	}

	<-p.exit
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stopping")
	close(p.exit)
	return nil
}
