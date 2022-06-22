package usecase

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/threads/domain"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
	"strconv"
)

type ThreadRepository interface {
	Create(ctx context.Context, thread domain.Thread) (domain.Thread, error)
	GetById(ctx context.Context, id int64) (domain.Thread, error)
	GetBySlug(ctx context.Context, slug string) (domain.Thread, error)
	Patch(ctx context.Context, id int64, partialThread domain.PartialThread) (domain.Thread, error)
	PatchBySlug(ctx context.Context, slug string, partialThread domain.PartialThread) (domain.Thread, error)
}

type ThreadUseCaseImpl struct {
	threadRepo ThreadRepository
}

func New(threadRepo ThreadRepository) *ThreadUseCaseImpl {
	return &ThreadUseCaseImpl{
		threadRepo: threadRepo,
	}
}

func (u *ThreadUseCaseImpl) Create(ctx context.Context, thread models.Thread) (models.Thread, error) {
	created, err := u.threadRepo.Create(ctx, domain.FromModel(thread, nil))
	return created.ToModel(), err
}

func (u *ThreadUseCaseImpl) GetById(ctx context.Context, id int64) (models.Thread, error) {
	obtained, err := u.threadRepo.GetById(ctx, id)
	return obtained.ToModel(), err
}

func (u *ThreadUseCaseImpl) GetBySlug(ctx context.Context, slug string) (models.Thread, error) {
	obtained, err := u.threadRepo.GetBySlug(ctx, slug)
	return obtained.ToModel(), err
}

func (u *ThreadUseCaseImpl) GetBySlugOrId(ctx context.Context, slugOrId string) (models.Thread, error) {
	id, err := strconv.ParseInt(slugOrId, 10, 64)

	if err != nil {
		obtained, err := u.threadRepo.GetById(ctx, id)
		return obtained.ToModel(), err
	} else {
		obtained, err := u.threadRepo.GetBySlug(ctx, slugOrId)
		return obtained.ToModel(), err
	}
}

func (u *ThreadUseCaseImpl) Patch(ctx context.Context, id int64, threadUpdate models.ThreadUpdate) (models.Thread, error) {
	edited, err := u.threadRepo.Patch(ctx, id, domain.FromModelUpdate(threadUpdate))
	return edited.ToModel(), err
}

func (u *ThreadUseCaseImpl) PatchBySlugOrid(ctx context.Context, slugOrId string, threadUpdate models.ThreadUpdate) (models.Thread, error) {
	id, err := strconv.ParseInt(slugOrId, 10, 64)

	if err != nil {
		edited, err := u.threadRepo.Patch(ctx, id, domain.FromModelUpdate(threadUpdate))
		return edited.ToModel(), err
	} else {
		edited, err := u.threadRepo.PatchBySlug(ctx, slugOrId, domain.FromModelUpdate(threadUpdate))
		return edited.ToModel(), err
	}
}
