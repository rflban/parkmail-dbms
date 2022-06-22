package usecase

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/votes/domain"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
	"strconv"
)

type VoteRepository interface {
	Set(ctx context.Context, vote domain.Vote) (domain.Vote, error)
	SetByThreadSlug(ctx context.Context, slug string, vote domain.Vote) (domain.Vote, error)
	Create(ctx context.Context, vote domain.Vote) (domain.Vote, error)
	Exists(ctx context.Context, nickname string, thread int64) (bool, error)
	Patch(ctx context.Context, nickname string, thread int64, voice *int64) (domain.Vote, error)
}

type VoteUseCaseImpl struct {
	voteRepo VoteRepository
}

func New(voteRepo VoteRepository) *VoteUseCaseImpl {
	return &VoteUseCaseImpl{
		voteRepo: voteRepo,
	}
}

func (u *VoteUseCaseImpl) Set(ctx context.Context, thread string, vote models.Vote) (models.Vote, error) {
	threadId, err := strconv.ParseInt(thread, 10, 64)
	toSet := domain.FromModel(vote, threadId)

	var (
		obtained domain.Vote
	)

	if err != nil {
		obtained, err = u.voteRepo.Set(ctx, toSet)
	} else {
		obtained, err = u.voteRepo.SetByThreadSlug(ctx, thread, toSet)
	}

	return obtained.ToModel(), err
}

func (u *VoteUseCaseImpl) Create(ctx context.Context, thread int64, vote models.Vote) (models.Vote, error) {
	created, err := u.voteRepo.Create(ctx, domain.FromModel(vote, thread))
	return created.ToModel(), err
}

func (u *VoteUseCaseImpl) Exists(ctx context.Context, nickname string, thread int64) (bool, error) {
	exists, err := u.voteRepo.Exists(ctx, nickname, thread)
	return exists, err
}

func (u *VoteUseCaseImpl) Patch(ctx context.Context, nickname string, thread int64, voice *int64) (models.Vote, error) {
	edited, err := u.voteRepo.Patch(ctx, nickname, thread, voice)
	return edited.ToModel(), err
}
