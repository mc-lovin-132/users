package service

import (
	"context"

	"github.com/mc-lovin-132/users/internal/domain"
)

type repository interface {
	Create(ctx context.Context, data *domain.User) (int, error)
	Update(ctx context.Context, id int, name, email, password *string) (int, error)
	Get(ctx context.Context, id *int, email *string) (*domain.User, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo repository
}

func New(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, data *domain.User) (int, error) {
	return s.repo.Create(ctx, data)
}
func (s *service) Update(ctx context.Context, id int, name, email, password *string) (int, error) {
	return s.repo.Update(ctx, id, name, email, password)
}
func (s *service) Get(ctx context.Context, id *int, email *string) (*domain.User, error) {
	if id != nil && email != nil {
		return nil, domain.ErrToMuchArgs
	} else if id == nil && email == nil {
		return nil, domain.ErrNotEnoughArgs
	}
	return s.repo.Get(ctx, id, email)
}
func (s *service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
