package main

import (
	"context"
	FasthttpRouter "github.com/fasthttp/router"
	"github.com/jackc/pgx/v4/pgxpool"
	ServiceDelivery "github.com/rflban/parkmail-dbms/internal/forum/service/delivery"
	ServiceRepo "github.com/rflban/parkmail-dbms/internal/forum/service/repository"
	ServiceUseCase "github.com/rflban/parkmail-dbms/internal/forum/service/usecase"
)

const prefix = "/api"

func SetupHandlers(ctx context.Context, pool *pgxpool.Pool, router *FasthttpRouter.Router) {
	serviceRepo := ServiceRepo.New(pool)

	serviceUseCase := ServiceUseCase.New(serviceRepo)

	serviceDelivery := ServiceDelivery.New(serviceUseCase)

	router.GET(prefix+"/service/status", serviceDelivery.Status)
	router.POST(prefix+"/service/clear", serviceDelivery.Clear)
}
