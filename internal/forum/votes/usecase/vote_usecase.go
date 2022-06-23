package usecase

import (
	"context"
	threadsDomain "github.com/rflban/parkmail-dbms/internal/forum/threads/domain"
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

type ThreadRepository interface {
	GetById(ctx context.Context, id int64) (threadsDomain.Thread, error)
	GetBySlug(ctx context.Context, slug string) (threadsDomain.Thread, error)
}

type VoteUseCaseImpl struct {
	voteRepo   VoteRepository
	threadRepo ThreadRepository
}

func New(voteRepo VoteRepository, threadRepo ThreadRepository) *VoteUseCaseImpl {
	return &VoteUseCaseImpl{
		voteRepo:   voteRepo,
		threadRepo: threadRepo,
	}
}

func (u *VoteUseCaseImpl) Set(ctx context.Context, thread string, vote models.Vote) (models.Thread, error) {
	threadId, err := strconv.ParseInt(thread, 10, 64)
	toSet := domain.FromModel(vote, threadId)

	var (
		threadEntity threadsDomain.Thread
	)

	isThreadSlug := err != nil

	if isThreadSlug {
		_, err = u.voteRepo.SetByThreadSlug(ctx, thread, toSet)
	} else {
		toSet.Thread = threadId
		_, err = u.voteRepo.Set(ctx, toSet)
	}

	if err == nil {
		if isThreadSlug {
			threadEntity, err = u.threadRepo.GetBySlug(ctx, thread)
		} else {
			toSet.Thread = threadId
			threadEntity, err = u.threadRepo.GetById(ctx, threadId)
		}
	}

	return threadEntity.ToModel(), err
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
