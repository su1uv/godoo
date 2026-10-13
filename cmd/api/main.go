package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"sync"
	"time"
)

var (
	buildTime string
	version   string
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	logger *slog.Logger
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment development|staging|production")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GODOO_DB_DSN"), "Postgres dsn")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-cons", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-cons", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
	}

	err = app.serve()
	if err != nil {
		logger.Error("failed to serve", "error", err)
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
