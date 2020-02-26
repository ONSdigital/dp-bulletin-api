package main

import (
	"github.com/ONSdigital/dp-bulletin-api/service"
	"github.com/ONSdigital/log.go/log"
	"os"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	if err := service.Run(BuildTime,GitCommit,Version, os.Args); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

