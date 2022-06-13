package usecase

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/user"
	"github.com/rflban/parkmail-dbms/internal/forum/user/domain"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
)

type UserUseCaseImpl struct {
	userRepo user.UserRepository
}

func New(userRepo user.UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userRepo: userRepo,
	}
}

func (u *UserUseCaseImpl) Create(ctx context.Context, user models.User) (models.User, error) {
	created, err := u.userRepo.Create(ctx, domain.GetUserEntity(user))
	return created.ToModel(), err
}

func (u *UserUseCaseImpl) Patch(ctx context.Context, nickname string, partialUser models.UserUpdate) (models.User, error) {
	updated, err := u.userRepo.Patch(ctx, nickname, domain.GetPartial(partialUser))
	return updated.ToModel(), err
}

func (u *UserUseCaseImpl) GetByEmail(ctx context.Context, email string) (models.User, error) {
	obtained, err := u.userRepo.GetByEmail(ctx, email)
	return obtained.ToModel(), err
}

func (u *UserUseCaseImpl) GetByNickname(ctx context.Context, nickname string) (models.User, error) {
	obtained, err := u.userRepo.GetByNickname(ctx, nickname)
	return obtained.ToModel(), err
}

func (u *UserUseCaseImpl) GetByEmailOrNickname(ctx context.Context, email, nickname string) (models.Users, error) {
	obtained, err := u.userRepo.GetByEmailOrNickname(ctx, email, nickname)

	users := make([]models.User, 0, len(obtained))
	for _, item := range obtained {
		users = append(users, item.ToModel())
	}

	return users, err
}
