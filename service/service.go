package service

import (
	"context"
	"github.com/ONSdigital/dp-bulletin-api/api"
	"github.com/ONSdigital/dp-bulletin-api/config"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"time"
)

const serviceName = "dp-bulletin-api"

// run the application
func Run(buildTime, gitCommit, version string, args []string) error {
	log.Namespace = serviceName
	cfg, err := config.Get()
	ctx := context.Background()
	if err != nil {
		return errors.Wrap(err, "unable to retrieve service configuration")
	}

	log.Event(ctx, "got service configuration", log.Data{"config": cfg}, log.INFO)

	versionInfo, err := healthcheck.NewVersionInfo(
		buildTime,
		gitCommit,
		version,
	)

	r := mux.NewRouter()

	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	if err = registerCheckers(ctx, &hc); err != nil {
		return err
	}
	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	a := api.Init(ctx, r, cfg)

	hc.Start(ctx)

	s := server.New(cfg.BindAddr, r)
	s.HandleOSSignals = false

	log.Event(ctx, "Starting server", log.Data{"config": cfg}, log.INFO)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Event(ctx, "failed to start http listen and serve", log.Error(err), log.ERROR)
			return
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Event(ctx, "shutting service down gracefully", log.INFO)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Event(ctx, "failed to shutdown http server", log.Error(err), log.ERROR)
	}
	if err := a.Close(ctx); err != nil {
		log.Event(ctx, "failed to shutdown api", log.Error(err), log.ERROR)
	}
	return nil
}

func registerCheckers(ctx context.Context, hc *healthcheck.HealthCheck) (err error) {
	// TODO ADD HEALTH CHECKS HERE
	return
}
