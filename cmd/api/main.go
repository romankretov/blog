package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"

	"local/blog/internal/app"
	"local/blog/internal/config"
	"local/blog/internal/observability"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Error().Err(err).Msg("failed to load config")
		os.Exit(1)
	}
	logger := observability.NewLogger(cfg.Env)
	a := app.New(cfg, logger)
	if err := a.Run(context.Background()); err != nil {
		logger.Error().Err(err).Msg("app exited with error")
		os.Exit(1)
	}

}
