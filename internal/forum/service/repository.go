package service

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/service/domain"
)

type ServiceRepository interface {
	Status(ctx context.Context) (domain.Status, error)
	Clear(ctx context.Context) error
}
