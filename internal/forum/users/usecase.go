package users

import (
	"context"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
)

type UserUseCase interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	Patch(ctx context.Context, nickname string, partialUser models.UserUpdate) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetByNickname(ctx context.Context, nickname string) (models.User, error)
	GetByEmailOrNickname(ctx context.Context, email, nickname string) (models.Users, error)
}
