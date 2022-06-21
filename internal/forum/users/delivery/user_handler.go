package delivery

import (
	"context"
	"encoding/json"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type UserUseCase interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	Patch(ctx context.Context, nickname string, partialUser models.UserUpdate) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetByNickname(ctx context.Context, nickname string) (models.User, error)
	GetByEmailOrNickname(ctx context.Context, email, nickname string) (models.Users, error)
}

type UserHandler struct {
	userUseCase UserUseCase
}

func New(userUseCase UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) Create(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	nickname, ok := rctx.UserValue("nickname").(string)
	if !ok {
		log.Errorf("Can't parse nickname: %v", rctx.UserValue("nickname"))
		body, _ := json.Marshal(models.Error{
			Message: "invalid nickname",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	toCreate := models.User{
		Nickname: &nickname,
	}
	if err := json.Unmarshal(rctx.PostBody(), &toCreate); err != nil {
		log.Error(err.Error())

		body, _ := json.Marshal(models.Error{
			Message: "invalid body",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	created, err := h.userUseCase.Create(ctx, toCreate)
	if err != nil {
		if _, ok := err.(forumErrors.UniqueError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "user already exists",
			})

			rctx.SetStatusCode(fasthttp.StatusConflict)
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

	body, err := json.Marshal(created)
	if err != nil {
		log.Error(err.Error())

		body, _ := json.Marshal(models.Error{
			Message: "internal server error",
		})

		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		rctx.SetBody(body)
		return
	}

	rctx.SetStatusCode(fasthttp.StatusCreated)
	rctx.SetBody(body)
}

func (h *UserHandler) GetProfileByNickname(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	nickname, ok := rctx.UserValue("nickname").(string)
	if !ok {
		log.Errorf("Can't parse nickname: %v", nickname)
		body, _ := json.Marshal(models.Error{
			Message: "invalid nickname",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	obtained, err := h.userUseCase.GetByNickname(ctx, nickname)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "user not found",
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

func (h *UserHandler) EditProfileByNickname(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	nickname, ok := rctx.UserValue("nickname").(string)
	if !ok {
		log.Errorf("Can't parse nickname: %v", nickname)
		body, _ := json.Marshal(models.Error{
			Message: "invalid nickname",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	toEdit := models.UserUpdate{}
	if err := json.Unmarshal(rctx.PostBody(), &toEdit); err != nil {
		log.Error(err.Error())

		body, _ := json.Marshal(models.Error{
			Message: "invalid body",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	edited, err := h.userUseCase.Patch(ctx, nickname, toEdit)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "user not found",
			})

			rctx.SetStatusCode(fasthttp.StatusNotFound)
			rctx.SetBody(body)
			return
		}

		if _, ok := err.(forumErrors.UniqueError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "conflict with another user's data",
			})

			rctx.SetStatusCode(fasthttp.StatusConflict)
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

	body, err := json.Marshal(edited)
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
