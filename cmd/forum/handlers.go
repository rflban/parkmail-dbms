package main

import (
	"context"
	FasthttpRouter "github.com/fasthttp/router"
	"github.com/jackc/pgx/v4/pgxpool"
	ServiceDelivery "github.com/rflban/parkmail-dbms/internal/forum/service/delivery"
	ServiceRepo "github.com/rflban/parkmail-dbms/internal/forum/service/repository"
	ServiceUseCase "github.com/rflban/parkmail-dbms/internal/forum/service/usecase"
	UserDelivery "github.com/rflban/parkmail-dbms/internal/forum/users/delivery"
	UserRepo "github.com/rflban/parkmail-dbms/internal/forum/users/repository"
	UserUseCase "github.com/rflban/parkmail-dbms/internal/forum/users/usecase"
)

const prefix = "/api"

func SetupHandlers(ctx context.Context, pool *pgxpool.Pool, router *FasthttpRouter.Router) {
	serviceRepo := ServiceRepo.New(pool)
	userRepo := UserRepo.New(pool)

	serviceUseCase := ServiceUseCase.New(serviceRepo)
	userUseCase := UserUseCase.New(userRepo)

	serviceHandler := ServiceDelivery.New(serviceUseCase)
	userHandler := UserDelivery.New(userUseCase)

	router.GET(prefix+"/service/status", serviceHandler.Status)
	router.POST(prefix+"/service/clear", serviceHandler.Clear)

	router.POST(prefix+"/user/{nickname}/create", userHandler.Create)
	router.GET(prefix+"/user/{nickname}/profile", userHandler.GetProfileByNickname)
	router.POST(prefix+"/user/{nickname}/profile", userHandler.EditProfileByNickname)
}
