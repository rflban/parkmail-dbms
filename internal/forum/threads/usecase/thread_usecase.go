package usecase

import (
	"context"
	forumsDomain "github.com/rflban/parkmail-dbms/internal/forum/forums/domain"
	"github.com/rflban/parkmail-dbms/internal/forum/threads/domain"
	usersDomain "github.com/rflban/parkmail-dbms/internal/forum/users/domain"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
	"strconv"
)

type ThreadRepository interface {
	Create(ctx context.Context, thread domain.Thread) (domain.Thread, error)
	GetById(ctx context.Context, id int64) (domain.Thread, error)
	GetBySlug(ctx context.Context, slug string) (domain.Thread, error)
	Patch(ctx context.Context, id int64, partialThread domain.PartialThread) (domain.Thread, error)
	PatchBySlug(ctx context.Context, slug string, partialThread domain.PartialThread) (domain.Thread, error)
}

type ForumRepository interface {
	GetBySlug(ctx context.Context, slug string) (forumsDomain.Forum, error)
}

type UserRepository interface {
	GetByNickname(ctx context.Context, nickname string) (usersDomain.User, error)
}

type ThreadUseCaseImpl struct {
	threadRepo ThreadRepository
	forumRepo  ForumRepository
	userRepo   UserRepository
}

func New(threadRepo ThreadRepository, forumRepo ForumRepository, userRepo UserRepository) *ThreadUseCaseImpl {
	return &ThreadUseCaseImpl{
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
		userRepo:   userRepo,
	}
}

func (u *ThreadUseCaseImpl) Create(ctx context.Context, thread models.Thread) (models.Thread, error) {
	if thread.Slug != nil {
		obtained, err := u.threadRepo.GetBySlug(ctx, *thread.Slug)
		if err == nil {
			return obtained.ToModel(), forumErrors.NewUniqueError("threads", "slug")
		}
	}

	if thread.Forum == nil {
		return thread, forumErrors.NewEntityNotExistsError("forums")
	}

	forum, err := u.forumRepo.GetBySlug(ctx, *thread.Forum)
	if err != nil {
		return thread, err
	}
	user, err := u.userRepo.GetByNickname(ctx, thread.Author)
	if err != nil {
		return thread, err
	}

	created, err := u.threadRepo.Create(ctx, domain.FromModel(thread, nil))
	created.Forum = forum.Slug
	created.Author = user.Nickname

	return created.ToModel(), err
}

func (u *ThreadUseCaseImpl) GetById(ctx context.Context, id int64) (models.Thread, error) {
	obtained, err := u.threadRepo.GetById(ctx, id)
	return obtained.ToModel(), err
}

func (u *ThreadUseCaseImpl) GetBySlug(ctx context.Context, slug string) (models.Thread, error) {
	obtained, err := u.threadRepo.GetBySlug(ctx, slug)
	return obtained.ToModel(), err
}

func (u *ThreadUseCaseImpl) GetBySlugOrId(ctx context.Context, slugOrId string) (models.Thread, error) {
	id, err := strconv.ParseInt(slugOrId, 10, 64)

	if err != nil {
		obtained, err := u.threadRepo.GetBySlug(ctx, slugOrId)
		return obtained.ToModel(), err
	} else {
		obtained, err := u.threadRepo.GetById(ctx, id)
		return obtained.ToModel(), err
	}
}

func (u *ThreadUseCaseImpl) Patch(ctx context.Context, id int64, threadUpdate models.ThreadUpdate) (models.Thread, error) {
	edited, err := u.threadRepo.Patch(ctx, id, domain.FromModelUpdate(threadUpdate))
	return edited.ToModel(), err
}

func (u *ThreadUseCaseImpl) PatchBySlugOrId(ctx context.Context, slugOrId string, threadUpdate models.ThreadUpdate) (models.Thread, error) {
	id, err := strconv.ParseInt(slugOrId, 10, 64)

	if err != nil {
		edited, err := u.threadRepo.Patch(ctx, id, domain.FromModelUpdate(threadUpdate))
		return edited.ToModel(), err
	} else {
		edited, err := u.threadRepo.PatchBySlug(ctx, slugOrId, domain.FromModelUpdate(threadUpdate))
		return edited.ToModel(), err
	}
}
