package server

import (
	"context"
	"errors"
	"net/http"

	"canary/internal/api"
	"canary/internal/config"
	"canary/internal/logger"
	serverpkg "canary/internal/server"
)

func Start() error {
	cfg := config.DefaultConfig()

	logger.Initialize(cfg.Logging.Level, cfg.Logging.Format)

	srv := serverpkg.New(":"+cfg.Server.Port, api.NewRouter(), logger.Get())

	errChan := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	return <-errChan
}

func Stop(ctx context.Context) error {
	return nil
}
