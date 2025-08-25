package storage

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/config"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/logger"
)

var DB *sql.DB

func InitDB(cfg *config.Config) {
	logger.Log.Info("starting db connection")
	dbParams := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)
	db, err := sql.Open("postgres", dbParams)
	if err != nil {
		logger.Log.Error("could not connect to db", slog.Any("err", err))
		panic("could not connect to db")
	}
	runMigrations(cfg)
	DB = db
	logger.Log.Info("db connected successfully")
}

func runMigrations(cfg *config.Config) {
	logger.Log.Info("starting migrations")
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		logger.Log.Error("migration setup error", slog.Any("err", err))
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.Error("migration failed", slog.Any("err", err))
	}
	logger.Log.Info("migrations applied successfully")
}
