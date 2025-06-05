package server

import (
	"canary/internal/config"
	"canary/internal/database"
	"context"
	"errors"
	"net/http"

	"canary/internal/api"
	"canary/internal/logger"
	serverpkg "canary/internal/server"
)

func Start() error {
	cfg := config.DefaultConfig()

	db, err := database.Connect(config.Config)
	if err != nil {
		return err
	}
	logger.Initialize(cfg.Logging.Level, cfg.Logging.Format)

	srv := serverpkg.New(":"+cfg.Server.Port, api.NewRouter(db), logger.Get())

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
