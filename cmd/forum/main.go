package main

import (
	"context"
	"fmt"
	FasthttpRouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	ctx := BindLoggers(context.Background())

	conf, err := getConfig(ctx)
	if err != nil {
		return
	}

	connString := GetConnString(conf.Database)
	pool, err := SetupDB(ctx, connString)
	if err != nil {
		return
	}

	router := FasthttpRouter.New()

	SetupHandlers(ctx, pool, router)

	fmt.Println(helloMessage)
	fmt.Printf("Server has been started at http://localhost:%d\n", conf.Server.Port)

	err = fasthttp.ListenAndServe(
		fmt.Sprintf(":%d", conf.Server.Port),
		func(fasthttpCtx *fasthttp.RequestCtx) {
			fasthttpCtx.SetUserValue("ctx", ctx)
			router.Handler(fasthttpCtx)
		},
	)

	if err != nil {
		fmt.Println(err)
	}
}
