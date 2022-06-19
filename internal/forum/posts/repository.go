package posts

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/posts/domain"
)

type PostRepository interface {
	CreateAt(ctx context.Context, forum string, thread int64, posts [][]interface{}) (domain.Post, error)
	Patch(ctx context.Context, id int64, message *string) (domain.Post, error)
	GetById(ctx context.Context, id int64) (domain.Post, error)
	GetDetails(ctx context.Context, id int64, related []string) (domain.PostFull, error)
	GetFromThreadFlat(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error)
	GetFromThreadTree(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error)
	GetFromThreadParentTree(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error)
}
