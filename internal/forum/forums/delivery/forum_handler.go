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

type ForumUseCase interface {
	Create(ctx context.Context, forum models.Forum) (models.Forum, error)
	GetBySlug(ctx context.Context, slug string) (models.Forum, error)
	GetUsersBySlug(ctx context.Context, slug string, since string, limit uint64, desc bool) (models.Users, error)
	GetThreadsBySlug(ctx context.Context, slug string, since string, limit uint64, desc bool) (models.Threads, error)
}

type ThreadUseCase interface {
	Create(ctx context.Context, thread models.Thread) (models.Thread, error)
}

type ForumHandler struct {
	forumUseCase  ForumUseCase
	threadUseCase ThreadUseCase
}

func New(forumUseCase ForumUseCase, threadUseCase ThreadUseCase) *ForumHandler {
	return &ForumHandler{
		forumUseCase:  forumUseCase,
		threadUseCase: threadUseCase,
	}
}

func (h *ForumHandler) Create(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	var fromBody models.Forum
	if err := json.Unmarshal(rctx.PostBody(), &fromBody); err != nil {
		log.Error(err.Error())

		body, _ := json.Marshal(models.Error{
			Message: "invalid body",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	obtained, err := h.forumUseCase.Create(ctx, fromBody)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "thread not found",
			})

			rctx.SetStatusCode(fasthttp.StatusNotFound)
			rctx.SetBody(body)
			return
		}

		if _, ok := err.(forumErrors.UniqueError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "conflict with another forum's data",
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

	rctx.SetStatusCode(fasthttp.StatusCreated)
	rctx.SetBody(body)
}

func (h *ForumHandler) GetDetails(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	slug, ok := rctx.UserValue("slug").(string)
	if !ok {
		log.Errorf("Can't parse slug: %v", rctx.UserValue("slug"))
		body, _ := json.Marshal(models.Error{
			Message: "invalid slug",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	obtained, err := h.forumUseCase.GetBySlug(ctx, slug)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "forum not found",
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

func (h *ForumHandler) CreateThread(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	slug, ok := rctx.UserValue("slug").(string)
	if !ok {
		log.Errorf("Can't parse slug: %v", rctx.UserValue("slug"))
		body, _ := json.Marshal(models.Error{
			Message: "invalid slug",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	var fromBody models.Thread
	if err := json.Unmarshal(rctx.PostBody(), &fromBody); err != nil {
		log.Error(err.Error())

		body, _ := json.Marshal(models.Error{
			Message: "invalid body",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	fromBody.Forum = &slug
	obtained, err := h.threadUseCase.Create(ctx, fromBody)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "thread or author not found",
			})

			rctx.SetStatusCode(fasthttp.StatusNotFound)
			rctx.SetBody(body)
			return
		}

		if _, ok := err.(forumErrors.UniqueError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "conflict with another thread's data",
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

	rctx.SetStatusCode(fasthttp.StatusCreated)
	rctx.SetBody(body)
}

func (h *ForumHandler) GetUsers(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	slug, ok := rctx.UserValue("slug").(string)
	if !ok {
		log.Errorf("Can't parse slug: %v", rctx.UserValue("slug"))
		body, _ := json.Marshal(models.Error{
			Message: "invalid slug",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	sinceRaw := rctx.QueryArgs().Peek("since")
	limitRaw := rctx.QueryArgs().Peek("limit")
	descRaw := rctx.QueryArgs().Peek("desc")

	since := string(sinceRaw)
	limit, _ := strconv.ParseUint(string(limitRaw), 10, 64)
	desc := string(descRaw) == "true"

	obtained, err := h.forumUseCase.GetUsersBySlug(ctx, slug, since, limit, desc)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "forum not found",
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

func (h *ForumHandler) GetThreads(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	slug, ok := rctx.UserValue("slug").(string)
	if !ok {
		log.Errorf("Can't parse slug: %v", rctx.UserValue("slug"))
		body, _ := json.Marshal(models.Error{
			Message: "invalid slug",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	sinceRaw := rctx.QueryArgs().Peek("since")
	limitRaw := rctx.QueryArgs().Peek("limit")
	descRaw := rctx.QueryArgs().Peek("desc")

	since := string(sinceRaw)
	limit, _ := strconv.ParseUint(string(limitRaw), 10, 64)
	desc := string(descRaw) == "true"

	obtained, err := h.forumUseCase.GetThreadsBySlug(ctx, slug, since, limit, desc)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "forum not found",
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
