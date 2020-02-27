package main

import (
	"github.com/ONSdigital/dp-bulletin-api/service"
	"github.com/ONSdigital/log.go/log"
	"os"
	"os/signal"
)

const serviceName = "dp-bulletin-api"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	log.Namespace = serviceName

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	svcErrors := make(chan error, 1)
	svc, err := service.Run(BuildTime, GitCommit, Version, svcErrors)
	if err != nil {
		log.Event(nil, "running service failed", log.Error(err), log.FATAL)
		os.Exit(1)
	}

	// blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		log.Event(nil, "service error received", log.Error(err), log.FATAL)
		os.Exit(1)
	case sig := <-signals:
		log.Event(nil, "os signal received", log.Data{"signal": sig}, log.INFO)
		svc.Close()
	}
}
