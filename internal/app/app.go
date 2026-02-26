package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"local/blog/internal/config"
	httpapi "local/blog/internal/http"

	"github.com/rs/zerolog"
)

type App struct {
	cfg config.Config
	log zerolog.Logger
}

func New(cfg config.Config, log zerolog.Logger) *App {
	return &App{cfg: cfg, log: log}
}

func (a *App) Run(ctx context.Context) error {
	deps, err := NewDeps(ctx, a.cfg.DatabaseURL, a.cfg.RedisAddr, a.cfg.RedisPassword, a.cfg.RedisDB)
	if err != nil {
		return err
	}
	defer deps.DB.Close()
	defer func() { _ = deps.Redis.Close() }()

	router := httpapi.NewRouter(a.log, deps.DB, deps.Redis)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.Port),
		Handler:      router,
		ReadTimeout:  a.cfg.ReadTimeout,
		WriteTimeout: a.cfg.WriteTimeout,
		IdleTimeout:  a.cfg.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		a.log.Info().Str("addr", srv.Addr).Msg("http server starting")
		errCh <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		a.log.Info().Msg("context canceled, shutting down")
	case sig := <-quit:
		a.log.Info().Str("signal", sig.String()).Msg("signal received, shutting down")
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.log.Info().Msg("http server stopping")
	return srv.Shutdown(shutdownCtx)
}
