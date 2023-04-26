package main

import (
	"flag"
	"log"

	"github.com/kardianos/service"
)

var (
	configFile *string
)

func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	configFile = flag.String("config", "/etc/fsmon.yaml", "Config file path.")
	flag.Parse()

	options := make(service.KeyValue)
	options["Restart"] = "always"

	svcConfig := &service.Config{
		Name:        "fsmon",
		DisplayName: "Filesystem monitoring",
		Description: "Filesystem monitoring tool.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if svcFlag != nil && *svcFlag != "" {
		err = service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
