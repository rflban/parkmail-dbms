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

type ThreadUseCase interface {
	Create(ctx context.Context, thread models.Thread) (models.Thread, error)
	GetBySlugOrId(ctx context.Context, slugOrId string) (models.Thread, error)
	PatchBySlugOrId(ctx context.Context, slugOrId string, threadUpdate models.ThreadUpdate) (models.Thread, error)
}

type PostUseCase interface {
	Create(ctx context.Context, threadSlugOrId string, posts models.Posts) (models.Posts, error)
	GetFromThread(ctx context.Context, thread string, since int64, limit uint64, desc bool, sort string) (models.Posts, error)
}

type VoteUseCase interface {
	Set(ctx context.Context, thread string, vote models.Vote) (models.Thread, error)
}

type ThreadHandler struct {
	threadUseCase ThreadUseCase
	postUseCase   PostUseCase
	voteUseCase   VoteUseCase
}

func New(postUseCase PostUseCase, threadUseCase ThreadUseCase, voteUseCase VoteUseCase) *ThreadHandler {
	return &ThreadHandler{
		threadUseCase: threadUseCase,
		postUseCase:   postUseCase,
		voteUseCase:   voteUseCase,
	}
}

func (h *ThreadHandler) CreatePosts(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	slugOrId, ok := rctx.UserValue("slug_or_id").(string)
	if !ok {
		log.Errorf("Can't parse slug: %v", rctx.UserValue("slug_or_id"))
		body, _ := json.Marshal(models.Error{
			Message: "invalid slug_or_id",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	var fromBody models.Posts
	if err := json.Unmarshal(rctx.PostBody(), &fromBody); err != nil {
		log.Error(err.Error())

		body, _ := json.Marshal(models.Error{
			Message: "invalid body",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	obtained, err := h.postUseCase.Create(ctx, slugOrId, fromBody)
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

func (h *ThreadHandler) GetDetails(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	slugOrId, ok := rctx.UserValue("slug_or_id").(string)
	if !ok {
		log.Errorf("Can't parse slug: %v", rctx.UserValue("slug_or_id"))
		body, _ := json.Marshal(models.Error{
			Message: "invalid slug_or_id",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	obtained, err := h.threadUseCase.GetBySlugOrId(ctx, slugOrId)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "thread not found",
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

func (h *ThreadHandler) Edit(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	slugOrId, ok := rctx.UserValue("slug_or_id").(string)
	if !ok {
		log.Errorf("Can't parse slug: %v", rctx.UserValue("slug_or_id"))
		body, _ := json.Marshal(models.Error{
			Message: "invalid slug_or_id",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	var fromBody models.ThreadUpdate
	if err := json.Unmarshal(rctx.PostBody(), &fromBody); err != nil {
		log.Error(err.Error())

		body, _ := json.Marshal(models.Error{
			Message: "invalid body",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	obtained, err := h.threadUseCase.PatchBySlugOrId(ctx, slugOrId, fromBody)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "thread not found",
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

func (h *ThreadHandler) GetPosts(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	slugOrId, ok := rctx.UserValue("slug_or_id").(string)
	if !ok {
		log.Errorf("Can't parse slug: %v", rctx.UserValue("slug_or_id"))
		body, _ := json.Marshal(models.Error{
			Message: "invalid slug_or_id",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	sortRaw := rctx.QueryArgs().Peek("sort")
	sinceRaw := rctx.QueryArgs().Peek("since")
	limitRaw := rctx.QueryArgs().Peek("limit")
	descRaw := rctx.QueryArgs().Peek("desc")

	sort := string(sortRaw)
	desc := string(descRaw) == "true"
	since, err := strconv.ParseInt(string(sinceRaw), 10, 64)
	if err != nil {
		since = 0
	}
	limit, err := strconv.ParseUint(string(limitRaw), 10, 64)
	if err != nil {
		limit = 0
	}

	obtained, err := h.postUseCase.GetFromThread(ctx, slugOrId, since, limit, desc, sort)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "thread not found",
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

func (h *ThreadHandler) Vote(rctx *fasthttp.RequestCtx) {
	ctx := rctx.UserValue("ctx").(context.Context)
	log := ctx.Value(constants.DeliveryLogKey).(*logrus.Entry)
	rctx.SetContentType("application/json")

	slugOrId, ok := rctx.UserValue("slug_or_id").(string)
	if !ok {
		log.Errorf("Can't parse slug: %v", rctx.UserValue("slug_or_id"))
		body, _ := json.Marshal(models.Error{
			Message: "invalid slug_or_id",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	var fromBody models.Vote
	if err := json.Unmarshal(rctx.PostBody(), &fromBody); err != nil {
		log.Error(err.Error())

		body, _ := json.Marshal(models.Error{
			Message: "invalid body",
		})

		rctx.SetStatusCode(fasthttp.StatusBadRequest)
		rctx.SetBody(body)
		return
	}

	obtained, err := h.voteUseCase.Set(ctx, slugOrId, fromBody)
	if err != nil {
		if _, ok := err.(forumErrors.EntityNotExistsError); ok {
			body, _ := json.Marshal(models.Error{
				Message: "thread or user not found",
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
