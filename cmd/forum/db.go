package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

func SetupDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
