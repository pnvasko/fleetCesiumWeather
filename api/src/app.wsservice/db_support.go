package app_wsservice

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgxpool"
)

func initDBDriver(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool

	dburl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?pool_max_conns=%d",
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.Database,
		cfg.Db.MaxConnections,
	)

	config, err := pgxpool.ParseConfig(fmt.Sprintf("%s&timezone=UTC", dburl))
	if err != nil {
		return nil, fmt.Errorf("error get DB config: %s", err)
	}

	config.MaxConns = int32(cfg.Db.MaxConnections)
	config.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		return nil
	}

	if pool, err = pgxpool.ConnectConfig(ctx, config); err != nil {
		return nil, fmt.Errorf("error connect to DB server: %s", err)
	}
	return pool, nil
}
