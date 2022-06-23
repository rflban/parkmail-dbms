package middlewares

import (
	"context"
	"github.com/google/uuid"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func AccessLog(next func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx) {
	return func(rctx *fasthttp.RequestCtx) {
		ctx := rctx.UserValue("ctx").(context.Context)
		log := ctx.Value(constants.AccessLogKey).(*logrus.Entry)

		connUuid := uuid.New()

		log.
			WithField("conn_uuid", connUuid).
			WithField("access", "request").
			WithField("body_size", len(rctx.Request.Body())).
			Info(string(rctx.Method()), " ", rctx.URI().String())
		next(rctx)
		log.
			WithField("conn_uuid", connUuid).
			WithField("access", "response").
			WithField("body_size", len(rctx.Response.Body())).
			Info("STATUS ", rctx.Response.StatusCode())
	}
}
