package forums

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/forums/domain"
	threadsDomain "github.com/rflban/parkmail-dbms/internal/forum/threads/domain"
	usersDomain "github.com/rflban/parkmail-dbms/internal/forum/users/domain"
)

type ForumRepository interface {
	Create(ctx context.Context, forum domain.Forum) (domain.Forum, error)
	GetBySlug(ctx context.Context, slug string) (domain.Forum, error)
	GetUsersBySlug(ctx context.Context, slug string, since string, limit uint64, desc bool) ([]usersDomain.User, error)
	GetThreadsBySlug(ctx context.Context, slug string, since string, limit uint64, desc bool) ([]threadsDomain.Thread, error)
}
