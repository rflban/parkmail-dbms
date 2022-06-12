package service

import (
	"context"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
)

type ServiceUseCase interface {
	Status(ctx context.Context) (models.Status, error)
	Clear(ctx context.Context) error
}
