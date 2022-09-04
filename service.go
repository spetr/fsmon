package main

import (
	"net/http"

	"github.com/kardianos/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	r := prometheus.NewRegistry()

	if conf.Prometheus != "" {
		http.Handle("/metrics", promhttp.HandlerFor(
			r,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		))
		go func() {
			if err := http.ListenAndServe(conf.Prometheus, nil); err != nil {
				logger.Error(err.Error())
			}
		}()

	}

	go monFsUpdate(r)

	<-p.exit
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stopping")
	close(p.exit)
	return nil
}
