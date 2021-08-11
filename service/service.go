package service

import (
	"context"
	"github.com/ONSdigital/dp-bulletin-api/api"
	"github.com/ONSdigital/dp-bulletin-api/config"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Service struct {
	Config      *config.Config
	Server      *server.Server
	Router      *mux.Router
	API         *api.API
	HealthCheck *healthcheck.HealthCheck
}

// Run the service
func Run(buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {
	ctx := context.Background()
	log.Info(ctx, "running service")

	cfg, err := config.Get()
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve service configuration")
	}
	log.Info(ctx, "got service configuration", log.Data{"config": cfg})

	r := mux.NewRouter()

	s := server.New(cfg.BindAddr, r)
	s.HandleOSSignals = false

	a := api.Init(ctx, r)

	versionInfo, err := healthcheck.NewVersionInfo(
		buildTime,
		gitCommit,
		version,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse version information")
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	if err := registerCheckers(ctx, &hc); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}
	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	hc.Start(ctx)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in  http listen and serve")
		}
	}()

	return &Service{
		Config:      cfg,
		Router:      r,
		API:         a,
		HealthCheck: &hc,
		Server:      s,
	}, nil
}

// Gracefully shutdown the service
func (svc *Service) Close() {
	timeout := svc.Config.GracefulShutdownTimeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})

	// stop any incoming requests before closing any outbound connections
	if err := svc.Server.Shutdown(ctx); err != nil {
		log.Error(ctx, "failed to shutdown http server", err)
	}

	if err := svc.API.Close(ctx); err != nil {
		log.Error(ctx, "error closing API", err)
	}

	log.Info(ctx, "graceful shutdown complete")
}

func registerCheckers(ctx context.Context, hc *healthcheck.HealthCheck) (err error) {
	// TODO ADD HEALTH CHECKS HERE
	return
}
