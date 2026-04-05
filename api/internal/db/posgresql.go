package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pq "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PostgresqlDB struct {
	host     string
	user     string
	port     int
	password string
	database string
}

func NewPostgresqlDB(
	port int,
	host,
	user,
	pwd,
	dbName string,
) *PostgresqlDB {

	return &PostgresqlDB{
		host:     host,
		user:     user,
		port:     port,
		password: pwd,
		database: dbName,
	}

}

func (p *PostgresqlDB) Connect(lc fx.Lifecycle, log *zap.Logger) (*sql.DB, error) {
	cfg := pq.Config{
		Host:     p.host,
		User:     p.user,
		Port:     uint16(p.port),
		Password: p.password,
		Database: p.database,
	}

	connector, err := pq.NewConnectorConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("building connector: %w", err)
	}

	db := sql.OpenDB(connector)
	// TODO: set from config
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := db.PingContext(ctx); err != nil {
				return fmt.Errorf("database unreachable on start: %w", err)
			}

			log.Info("Database connected at", zap.String("addr", fmt.Sprintf("%v:%v", p.host, p.port)))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := db.Close(); err != nil {
				return fmt.Errorf("closing database: %w", err)
			}

			log.Info("Database closed at", zap.String("addr", fmt.Sprintf("%v:%v", p.host, p.port)))
			return nil
		},
	})

	return db, nil
}
