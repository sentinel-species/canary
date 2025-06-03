package server

import (
	"canary/internal/config"
	"context"
	"errors"
	"net/http"

	"canary/internal/api"
	"canary/internal/logger"
	serverpkg "canary/internal/server"
	db "canary/pkg/database"
)

func Start() error {
	cfg := config.DefaultConfig()

	_, err := db.Connect(config.Config)
	if err != nil {
		return err
	}

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
