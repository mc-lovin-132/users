package service

import (
	"context"

	"github.com/jmoiron/sqlx"
)

const (
	await = "await"
	alive = "alive"
	dead  = "dead"
)

type healthService struct {
	db *sqlx.DB
}

func NewHealthService(db *sqlx.DB) *healthService {
	return &healthService{db: db}
}

func (h *healthService) Status(ctx context.Context) (string, error) {
	if h.db == nil {
		return await, nil
	}
	err := h.db.PingContext(ctx)
	if err != nil {
		return dead, err
	}
	return alive, nil
}
