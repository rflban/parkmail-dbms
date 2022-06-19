package posts

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/posts/domain"
)

type PostRepository interface {
	Create(ctx context.Context, posts []domain.Post) ([]domain.Post, error)
	Patch(ctx context.Context, id int64, message *string) (domain.Post, error)
	GetById(ctx context.Context, id int64) (domain.Post, error)
	GetFromThreadFlat(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error)
	GetFromThreadTree(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error)
	GetFromThreadParentTree(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error)
}
