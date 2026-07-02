package repository

import (
	"context"
	"database/sql"

	"github.com/mc-lovin-132/users/internal/domain"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

// добавить маппинг ошибок дб в доменные
func (r *repository) Create(ctx context.Context, data *domain.User) (int, error) {
	var id int
	model := fromDomain(data)
	rows, err := r.db.NamedQueryContext(ctx, createQuery(), model)
	if err != nil {
		return 0, errorMapper(err)
	}
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return 0, errorMapper(err)
		}
	} else {
		return 0, errorMapper(sql.ErrNoRows)
	}
	return id, nil
}
func (r *repository) Update(ctx context.Context, id int, name, email, password *string) (int, error) {
	query, args, err := updateQuery(id, name, email, password)
	if err != nil {
		return 0, errorMapper(err)
	}
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, errorMapper(err)
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return 0, errorMapper(sql.ErrNoRows)
	}
	return id, nil
}
func (r *repository) Get(ctx context.Context, id *int, email *string) (*domain.User, error) {
	var model UserModel
	query, args, err := getQuery(id, email)
	if err != nil {
		return nil, errorMapper(err)
	}
	err = r.db.GetContext(ctx, &model, query, args...)
	if err != nil {
		return nil, errorMapper(err)
	}
	return toDomain(&model), nil
}
func (r *repository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, deleteQuery(), id)
	if err != nil {
		return errorMapper(err)
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errorMapper(sql.ErrNoRows)
	}
	return nil
}
