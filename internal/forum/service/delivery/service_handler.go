package delivery

import (
	"context"
	"encoding/json"
	"github.com/rflban/parkmail-dbms/internal/forum/service"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type ServiceHandler struct {
	serviceUseCase service.ServiceUseCase
}

func New(serviceUseCase service.ServiceUseCase) *ServiceHandler {
	return &ServiceHandler{
		serviceUseCase: serviceUseCase,
	}
}

func (h *ServiceHandler) Status(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)

	status, err := h.serviceUseCase.Status(ctx)
	if err != nil {
		body, _ := json.Marshal(models.Error{
			Message: "internal server error",
		})

		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		rctx.SetBody(body)
		return
	}

	body, err := json.Marshal(status)
	if err != nil {
		log.Error(err)

		body, _ := json.Marshal(models.Error{
			Message: "internal server error",
		})

		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		rctx.SetBody(body)
		return
	}

	rctx.SetStatusCode(fasthttp.StatusOK)
	rctx.SetContentType("application/json")
	rctx.SetBody(body)
}

func (h *ServiceHandler) Clear(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)

	err := h.serviceUseCase.Clear(ctx)
	if err != nil {
		body, _ := json.Marshal(models.Error{
			Message: "internal server error",
		})

		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		rctx.SetBody(body)
		return
	}

	rctx.SetStatusCode(fasthttp.StatusOK)
}
