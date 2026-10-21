package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/su1uv/godoo/internal/db/sqlc"
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
	routes  *http.ServeMux
	queries *db.Queries
	wg      sync.WaitGroup
}

func main() {
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
		os.Exit(1)
	}
	defer pool.Close()

	logger.Info("database connection pool established")

	queries := db.New(pool)
	routes := registerRoutes()

	app := &application{
		config:  cfg,
		logger:  logger,
		queries: queries,
		routes:  routes,
	}

	err = app.serve()
	if err != nil {
		logger.Error("failed to serve", "error", err)
		os.Exit(1)
	}
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
