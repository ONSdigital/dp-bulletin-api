package main

import (
	"context"
	"github.com/ONSdigital/dp-bulletin-api/api"
	"github.com/ONSdigital/dp-bulletin-api/config"
	"os"
	"os/signal"
	"time"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
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
	if err := run(os.Args); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

// run the application
func run(args []string) error {
	log.Namespace = "dp_bulletin_api"
	cfg, err := config.Get()
	ctx := context.Background()
	if err != nil {
		log.Event(ctx, "unable to retrieve service configuration", log.Error(err))
		return err
	}

	log.Event(ctx, "got service configuration", log.Data{"config": cfg})

	versionInfo, err := health.NewVersionInfo(
		BuildTime,
		GitCommit,
		Version,
	)

	r := mux.NewRouter()

	healthcheck := health.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	if err = registerCheckers(ctx, &healthcheck); err != nil {
		return err
	}
	r.StrictSlash(true).Path("/health").HandlerFunc(healthcheck.Handler)
	a := api.Init(ctx, r, cfg)

	healthcheck.Start(ctx)

	s := server.New(cfg.BindAddr, r)
	s.HandleOSSignals = false

	log.Event(ctx, "Starting server", log.Data{"config": cfg})

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Event(ctx, "failed to start http listen and serve", log.Error(err))
			return
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Event(ctx, "shutting service down gracefully")
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Event(ctx, "failed to shutdown http server", log.Error(err))
	}
	if err := a.Close(ctx); err != nil {
		log.Event(ctx, "failed to shutdown api", log.Error(err))
	}
	return nil
}

func registerCheckers(ctx context.Context, h *health.HealthCheck) (err error) {
	// TODO ADD HEALTH CHECKS HERE
	return
}
