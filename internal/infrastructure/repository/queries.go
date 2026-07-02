package repository

import (
	"github.com/mc-lovin-132/users/internal/domain"

	"github.com/Masterminds/squirrel"
)

// INSERT INTO
// users (name, email, password)
// VALUES (:name, :email, :password)
// RETURNING id;
func createQuery() string {
	return `
	INSERT INTO
	users (name, email, password)
	VALUES (:name, :email, :password)
	RETURNING id;`
}

// UPDATE users SET
//
//	name = $1,
//	email = $2,
//	password = $3
//
// WHERE id = $4
func updateQuery(id int, name, email, password *string) (string, []interface{}, error) {
	if name == nil && email == nil && password == nil {
		// TODO: возможно стоит заменить ошибку
		return "", nil, domain.ErrNotEnoughArgs
	}
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	builder := psql.Update("users")
	if name != nil {
		builder = builder.Set("name", *name)
	}
	if email != nil {
		builder = builder.Set("email", *email)
	}
	if password != nil {
		builder = builder.Set("password", *password)
	}

	return builder.Where(squirrel.Eq{"id": id}).ToSql()
}

// SELECT * FROM users WHERE id = :id / email
func getQuery(id *int, email *string) (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	builder := psql.Select("*").From("users")
	if id != nil {
		builder = builder.Where(squirrel.Eq{"id": *id})
	} else {
		builder = builder.Where(squirrel.Eq{"email": *email})
	}
	return builder.ToSql()
}

// DELETE FROM users WHERE id = :id
func deleteQuery() string {
	return `DELETE FROM users WHERE id = $1`
}
