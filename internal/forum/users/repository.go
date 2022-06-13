package users

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/users/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	Patch(ctx context.Context, nickname string, partialUser domain.PartialUser) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByNickname(ctx context.Context, nickname string) (domain.User, error)
	GetByEmailOrNickname(ctx context.Context, email, nickname string) ([]domain.User, error)
}
