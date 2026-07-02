package repository

import (
	"database/sql"

	"github.com/mc-lovin-132/users/internal/domain"

	"github.com/lib/pq"
)

type UserModel struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

// TODO: add error mapper
func errorMapper(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		// нарушение уникальности
		if pqErr.Code == "23505" {
			return domain.ErrNotUniqueEmail
		}
	}
	if err == sql.ErrNoRows {
		return domain.ErrUserNotFound
	}
	return domain.ErrInternal
}

func fromDomain(data *domain.User) *UserModel {
	return &UserModel{
		ID:       data.ID,
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	}
}

func toDomain(data *UserModel) *domain.User {
	return &domain.User{
		ID:       data.ID,
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	}
}
