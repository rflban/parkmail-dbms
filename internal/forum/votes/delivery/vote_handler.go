package delivery

import (
	"context"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
)

type VoteUseCase interface {
	Create(ctx context.Context, vote models.Vote) (models.Vote, error)
	Exists(ctx context.Context, nickname string, thread int64) (bool, error)
	Patch(ctx context.Context, nickname string, thread int64, voice *int64) (models.Vote, error)
}
