package delivery

import (
	"context"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
)

type ThreadUseCase interface {
	Create(ctx context.Context, thread models.Thread) (models.Thread, error)
	GetById(ctx context.Context, id int64) (models.Thread, error)
	GetBySlug(ctx context.Context, slug string) (models.Thread, error)
	Patch(ctx context.Context, id int64, threadUpdate models.ThreadUpdate) (models.Thread, error)
}
