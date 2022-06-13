package delivery

import (
	"context"
	"encoding/json"
	"github.com/rflban/parkmail-dbms/internal/forum/user"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	userUseCase user.UserUseCase
}

func New(userUseCase user.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) Create(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)

	nickname, ok := rctx.UserValue("nickname").(string)
	if !ok {
		log.Errorf("Can't parse nickname: %v", nickname)
		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	toCreate := models.User{
		Nickname: &nickname,
	}
	if err := json.Unmarshal(rctx.PostBody(), &toCreate); err != nil {
		log.Error(err.Error())
		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	created, err := h.userUseCase.Create(ctx, toCreate)
	if err != nil {
		if _, ok := err.(forumErrors.UniqueError); ok {
			rctx.SetStatusCode(fasthttp.StatusConflict)
			return
		}

		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(created)
	if err != nil {
		log.Error(err.Error())
		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	rctx.SetStatusCode(fasthttp.StatusCreated)
	rctx.SetContentType("application/json")
	rctx.SetBody(body)
}

func (h *UserHandler) GetProfileByNickname(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)

	nickname, ok := rctx.UserValue("nickname").(string)
	if !ok {
		log.Errorf("Can't parse nickname: %v", nickname)
		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	obtained, err := h.userUseCase.GetByNickname(ctx, nickname)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			rctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}

		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(obtained)
	if err != nil {
		log.Error(err.Error())
		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	rctx.SetStatusCode(fasthttp.StatusOK)
	rctx.SetContentType("application/json")
	rctx.SetBody(body)
}

func (h *UserHandler) EditProfileByNickname(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)

	nickname, ok := rctx.UserValue("nickname").(string)
	if !ok {
		log.Errorf("Can't parse nickname: %v", nickname)
		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	toEdit := models.UserUpdate{}
	if err := json.Unmarshal(rctx.PostBody(), &toEdit); err != nil {
		log.Error(err.Error())
		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	edited, err := h.userUseCase.Patch(ctx, nickname, toEdit)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			rctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	body, err := json.Marshal(edited)
	if err != nil {
		log.Error(err.Error())
		rctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	rctx.SetStatusCode(fasthttp.StatusOK)
	rctx.SetContentType("application/json")
	rctx.SetBody(body)
}
