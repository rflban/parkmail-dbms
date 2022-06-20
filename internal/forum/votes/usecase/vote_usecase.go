package usecase

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/votes/domain"
)

type VoteRepository interface {
	Create(ctx context.Context, vote domain.Vote) (domain.Vote, error)
	Exists(ctx context.Context, nickname string, thread int64) (bool, error)
	Patch(ctx context.Context, nickname string, thread int64, voice *int64) (domain.Vote, error)
}
