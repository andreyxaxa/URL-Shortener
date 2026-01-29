package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreyxaxa/URL-Shortener/config"
	"github.com/andreyxaxa/URL-Shortener/internal/controller/restapi"
	"github.com/andreyxaxa/URL-Shortener/internal/repo/cache"
	"github.com/andreyxaxa/URL-Shortener/internal/repo/persistent"
	"github.com/andreyxaxa/URL-Shortener/internal/usecase/link"
	"github.com/andreyxaxa/URL-Shortener/pkg/httpserver"
	"github.com/andreyxaxa/URL-Shortener/pkg/logger"
	"github.com/andreyxaxa/URL-Shortener/pkg/postgres"
	"github.com/andreyxaxa/URL-Shortener/pkg/redis"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// Postgres Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %v", err))
	}
	defer pg.Close()

	// Redis cache
	rd, err := redis.New(cfg.Redis.Addr, redis.DB(cfg.Redis.DB))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - redis.New: %v", err))
	}
	defer rd.Close()

	// Use-Case
	linkUseCase := link.New(
		persistent.New(pg),
		cache.New(rd),
		l,
	)

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	restapi.NewRouter(httpServer.App, linkUseCase, l, fmt.Sprintf("http://localhost:%s", cfg.HTTP.Port))

	// Start server
	httpServer.Start()

	// Waiting Signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %v", err))
	}

	err = httpServer.Shutdown()
}
