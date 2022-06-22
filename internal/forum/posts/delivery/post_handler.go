package delivery

import (
	"context"
	"encoding/json"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"strconv"
)

type PostUseCase interface {
	Patch(ctx context.Context, id int64, message *string) (models.Post, error)
	GetDetails(ctx context.Context, id int64, related []string) (models.PostFull, error)
}

type PostHandler struct {
	postUseCase PostUseCase
}

func New(postUseCase PostUseCase) *PostHandler {
	return &PostHandler{
		postUseCase: postUseCase,
	}
}

func (h *PostHandler) GetDetails(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	var (
		id  int64
		err error
	)

	idRaw, ok := rctx.UserValue("id").(string)
	if ok {
		id, err = strconv.ParseInt(idRaw, 10, 64)
	}

	if !ok || err != nil {
		log.Errorf("Can't parse id: %v", rctx.UserValue("id"))
		if err != nil {
			log.Error(err.Error())
		}

		body, _ := json.Marshal(models.Error{
			Message: "invalid id",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	relatedRaw := rctx.QueryArgs().PeekMulti("related")
	related := make([]string, 0, len(relatedRaw))

	for _, entity := range relatedRaw {
		related = append(related, string(entity))
	}

	obtained, err := h.postUseCase.GetDetails(ctx, id, related)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "post not found",
			})

			rctx.SetStatusCode(fasthttp.StatusNotFound)
			rctx.SetBody(body)
			return
		}

		body, _ := json.Marshal(models.Error{
			Message: "internal server error",
		})

		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		rctx.SetBody(body)
		return
	}

	body, err := json.Marshal(obtained)
	if err != nil {
		log.Error(err.Error())

		body, _ := json.Marshal(models.Error{
			Message: "internal server error",
		})

		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		rctx.SetBody(body)
		return
	}

	rctx.SetStatusCode(fasthttp.StatusOK)
	rctx.SetBody(body)
}

func (h *PostHandler) Edit(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	var fromBody models.PostUpdate
	if err := json.Unmarshal(rctx.PostBody(), &fromBody); err != nil {
		log.Error(err.Error())

		body, _ := json.Marshal(models.Error{
			Message: "invalid body",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	var (
		id  int64
		err error
	)

	idRaw, ok := rctx.UserValue("id").(string)
	if ok {
		id, err = strconv.ParseInt(idRaw, 10, 64)
	}

	if !ok || err != nil {
		log.Errorf("Can't parse id: %v", rctx.UserValue("id"))
		if err != nil {
			log.Error(err.Error())
		}

		body, _ := json.Marshal(models.Error{
			Message: "invalid id",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	obtained, err := h.postUseCase.Patch(ctx, id, fromBody.Message)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "post not found",
			})

			rctx.SetStatusCode(fasthttp.StatusNotFound)
			rctx.SetBody(body)
			return
		}

		body, _ := json.Marshal(models.Error{
			Message: "internal server error",
		})

		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		rctx.SetBody(body)
		return
	}

	body, err := json.Marshal(obtained)
	if err != nil {
		log.Error(err.Error())

		body, _ := json.Marshal(models.Error{
			Message: "internal server error",
		})

		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		rctx.SetBody(body)
		return
	}

	rctx.SetStatusCode(fasthttp.StatusOK)
	rctx.SetBody(body)
}
