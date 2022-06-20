package usecase

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/forum/service/domain"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
)

type ServiceRepository interface {
	Status(ctx context.Context) (domain.Status, error)
	Clear(ctx context.Context) error
}

type ServiceUseCaseImpl struct {
	serviceRepo ServiceRepository
}

func New(serviceRepo ServiceRepository) *ServiceUseCaseImpl {
	return &ServiceUseCaseImpl{
		serviceRepo: serviceRepo,
	}
}

func (uc *ServiceUseCaseImpl) Status(ctx context.Context) (models.Status, error) {
	status, err := uc.serviceRepo.Status(ctx)
	return status.ToModel(), err
}

func (uc *ServiceUseCaseImpl) Clear(ctx context.Context) error {
	err := uc.serviceRepo.Clear(ctx)
	return err
}
