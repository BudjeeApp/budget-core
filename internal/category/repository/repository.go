package repository

import (
	"context"
	"github.com/BudjeeApp/budget-core/config"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Repository interface {
	CreateCategory(ctx context.Context, request Category) (*Category, error)
	GetCategory(ctx context.Context, id string) (*Category, error)
}

type repository struct {
	Pool   sqlx.DB
	Logger *zap.Logger
}

func NewRepository(db sqlx.DB) Repository {
	return repository{
		Pool:   db,
		Logger: config.Logger,
	}
}
