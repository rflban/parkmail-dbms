package delivery

import (
	"encoding/json"
	"github.com/rflban/parkmail-dbms/internal/forum/service"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type ServiceHandler struct {
	logger         *logrus.Entry
	serviceUseCase service.ServiceUseCase
}

func NewServiceHandler(logger *logrus.Entry, serviceUseCase service.ServiceUseCase) ServiceHandler {
	return ServiceHandler{
		logger:         logger,
		serviceUseCase: serviceUseCase,
	}
}

func (h *ServiceHandler) Status(ctx *fasthttp.RequestCtx) {
	status, err := h.serviceUseCase.Status(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(status)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(body)
}

func (h *ServiceHandler) Clear(ctx *fasthttp.RequestCtx) {
	err := h.serviceUseCase.Clear(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}
