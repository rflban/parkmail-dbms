package usecase

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/threads/domain"
)

type ThreadRepository interface {
	Create(ctx context.Context, thread domain.Thread) (domain.Thread, error)
	GetById(ctx context.Context, id int64) (domain.Thread, error)
	GetBySlug(ctx context.Context, slug string) (domain.Thread, error)
	Patch(ctx context.Context, id int64, partialThread domain.PartialThread) (domain.Thread, error)
}
