package app

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Deps struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
}

func NewDeps(ctx context.Context, dbUrl, redisAddr, redisPassword string, redisDB int) (*Deps, error) {
	poolCfg, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, err
	}
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	cctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := db.Ping(cctx); err != nil {
		db.Close()
		return nil, err
	}
	if err := rdb.Ping(cctx).Err(); err != nil {
		db.Close()
		_ = rdb.Close()
		return nil, err
	}
	_, err = db.Exec(ctx, "SELECT 1 FROM users LIMIT 1;")
	if err != nil {
		db.Close()
		rdb.Close()
		return nil, err
	}

	return &Deps{
		DB:    db,
		Redis: rdb,
	}, nil
}
