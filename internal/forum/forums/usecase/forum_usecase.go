package usecase

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/forums/domain"
	threadsDomain "github.com/rflban/parkmail-dbms/internal/forum/threads/domain"
	usersDomain "github.com/rflban/parkmail-dbms/internal/forum/users/domain"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
)

type ForumRepository interface {
	Create(ctx context.Context, forum domain.Forum) (domain.Forum, error)
	GetBySlug(ctx context.Context, slug string) (domain.Forum, error)
	GetUsersBySlug(ctx context.Context, slug string, since string, limit uint64, desc bool) ([]usersDomain.User, error)
	GetThreadsBySlug(ctx context.Context, slug string, since string, limit uint64, desc bool) ([]threadsDomain.Thread, error)
}

type ForumUseCaseImpl struct {
	forumRepo ForumRepository
}

func New(forumRepo ForumRepository) *ForumUseCaseImpl {
	return &ForumUseCaseImpl{
		forumRepo: forumRepo,
	}
}

func (u *ForumUseCaseImpl) Create(ctx context.Context, forum models.Forum) (models.Forum, error) {
	created, err := u.forumRepo.Create(ctx, domain.FromModel(forum, nil))

	if err == nil {
		return created.ToModel(), err
	}

	if _, isConflict := err.(forumErrors.UniqueError); !isConflict {
		return models.Forum{}, err
	}
	conflict := err

	existing, err := u.forumRepo.GetBySlug(ctx, forum.Slug)
	if err != nil {
		return models.Forum{}, err
	}

	return existing.ToModel(), conflict
}

func (u *ForumUseCaseImpl) GetBySlug(ctx context.Context, slug string) (models.Forum, error) {
	obtained, err := u.forumRepo.GetBySlug(ctx, slug)
	return obtained.ToModel(), err
}

func (u *ForumUseCaseImpl) GetUsersBySlug(ctx context.Context, slug string, since string, limit uint64, desc bool) (models.Users, error) {
	_, err := u.forumRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	obtained, err := u.forumRepo.GetUsersBySlug(ctx, slug, since, limit, desc)

	if err != nil {
		return nil, err
	}

	users := make(models.Users, 0, len(obtained))
	for _, user := range obtained {
		users = append(users, user.ToModel())
	}

	return users, err
}

func (u *ForumUseCaseImpl) GetThreadsBySlug(ctx context.Context, slug string, since string, limit uint64, desc bool) (models.Threads, error) {
	_, err := u.forumRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	obtained, err := u.forumRepo.GetThreadsBySlug(ctx, slug, since, limit, desc)

	if err != nil {
		return nil, err
	}

	threads := make(models.Threads, 0, len(obtained))
	for _, thread := range obtained {
		threads = append(threads, thread.ToModel())
	}

	return threads, err
}
