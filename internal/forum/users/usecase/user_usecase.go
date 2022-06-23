package usecase

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/users/domain"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	Patch(ctx context.Context, nickname string, partialUser domain.PartialUser) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByNickname(ctx context.Context, nickname string) (domain.User, error)
	GetByEmailOrNickname(ctx context.Context, email, nickname string) ([]domain.User, error)
}

type UserUseCaseImpl struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userRepo: userRepo,
	}
}

func (u *UserUseCaseImpl) Create(ctx context.Context, user models.User) (models.Users, error) {
	created, err := u.userRepo.Create(ctx, domain.GetUserEntity(user))

	if err == nil {
		return models.Users{created.ToModel()}, nil
	}

	if _, isConflict := err.(forumErrors.UniqueError); !isConflict {
		return nil, err
	}
	conflict := err

	var nickname string
	if user.Nickname != nil {
		nickname = *user.Nickname
	}

	existing, err := u.userRepo.GetByEmailOrNickname(ctx, user.Email, nickname)
	if err != nil {
		return nil, err
	}

	users := make(models.Users, 0, len(existing))
	for _, user := range existing {
		users = append(users, user.ToModel())
	}

	return users, conflict
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
