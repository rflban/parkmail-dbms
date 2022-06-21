package delivery

import (
	"context"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
)

type PostUseCase interface {
	Create(ctx context.Context, posts models.Posts) (models.Posts, error)
	Patch(ctx context.Context, id int64, message *string) (models.Post, error)
	GetById(ctx context.Context, id int64) (models.Post, error)
	GetDetails(ctx context.Context, id int64, related []string) (models.PostFull, error)
	GetFromThread(ctx context.Context, thread int64, since int64, limit uint64, desc bool, sort string) (models.Posts, error)
}
