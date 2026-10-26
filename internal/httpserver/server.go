package httpserver

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/su1uv/godoo/internal/platform/db/sqlc"
)

var (
	buildTime string
	version   string
)

type config struct {
	port int
	env  string
	db   struct {
		dsn               string
		maxConns          int
		minConns          int
		maxConnIdleTime   string
		maxConnLifetime   string
		healthCheckPeriod string
	}
}

type application struct {
	config  config
	logger  *slog.Logger
	queries *db.Queries
	wg      sync.WaitGroup
}

func Run() error {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment development|staging|production")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GODOO_DB_DSN"), "PostgreSQL dsn")
	flag.IntVar(&cfg.db.maxConns, "db-max-cons", 25, "PostgreSQL max connections")
	flag.IntVar(&cfg.db.minConns, "db-min-cons", 5, "PostgreSQL min connections")
	flag.StringVar(&cfg.db.maxConnIdleTime, "db-max-cons-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.StringVar(&cfg.db.maxConnLifetime, "db-max-cons-lifetime", "1h", "PostgreSQL max connection lifetime")
	flag.StringVar(&cfg.db.healthCheckPeriod, "db-health-check-period", "30s", "PostgreSQL health check period")

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := openDB(ctx, cfg)
	if err != nil {
		logger.Error("failed to open database", "error", err)
		return err
	}
	defer pool.Close()

	logger.Info("database connection pool established")

	queries := db.New(pool)

	app := &application{
		config:  cfg,
		logger:  logger,
		queries: queries,
	}

	err = app.serve()
	if err != nil {
		logger.Error("failed to serve", "error", err)
	}

	return err
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.Info("shutting down server",
			"env", app.config.env,
			"signal", s.String(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks",
			"env", app.config.env,
			"addr", srv.Addr,
		)

		app.wg.Wait()
		shutdownError <- nil

	}()

	app.logger.Info("Starting server",
		"env", app.config.env,
		"addr", srv.Addr,
	)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server",
		"env", app.config.env,
		"addr", srv.Addr,
	)

	return nil
}

func openDB(ctx context.Context, cfg config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	fmt.Printf("dsn: %q\n", cfg.db.dsn)

	idleDuration, err := time.ParseDuration(cfg.db.maxConnIdleTime)
	if err != nil {
		return nil, err
	}

	lifetimeDuration, err := time.ParseDuration(cfg.db.maxConnLifetime)
	if err != nil {
		return nil, err
	}

	healthCheckDuration, err := time.ParseDuration(cfg.db.healthCheckPeriod)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(cfg.db.maxConns)
	config.MinConns = int32(cfg.db.minConns)
	config.MaxConnIdleTime = idleDuration
	config.MaxConnLifetime = lifetimeDuration
	config.HealthCheckPeriod = healthCheckDuration

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
