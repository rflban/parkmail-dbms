package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	"github.com/sirupsen/logrus"
	"time"
)

type DBConnConfig struct {
	Name          string
	Username      string
	Password      string
	Host          string
	Port          int
	MaxConns      int
	MinConns      int
	MaxIdleTimeNS time.Duration
}

func GetConnString(args DBConnConfig) string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d pool_max_conns=%d pool_min_conns=%d pool_max_conn_idle_time=%dns",
		args.Username, args.Password, args.Name, args.Host, args.Port, args.MaxConns, args.MinConns, args.MaxIdleTimeNS,
	)
}

func SetupDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	log, hasLogger := ctx.Value(constants.SetupLogKey).(*logrus.Entry)

	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		if hasLogger {
			log.Error(err.Error())
		}
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		if hasLogger {
			log.Error(err.Error())
		}
		return nil, err
	}

	return pool, nil
}
