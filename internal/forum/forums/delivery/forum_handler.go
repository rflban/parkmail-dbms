package delivery

import (
	"context"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
)

type ForumUseCase interface {
	Create(ctx context.Context, forum models.Forum) (models.Forum, error)
	GetBySlug(ctx context.Context, slug string) (models.Forum, error)
	GetUsersBySlug(ctx context.Context, slug string, since string, limit uint64, desc bool) (models.Users, error)
	GetThreadsBySlug(ctx context.Context, slug string, since string, limit uint64, desc bool) (models.Threads, error)
}
