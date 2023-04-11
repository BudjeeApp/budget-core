package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateUser(ctx context.Context, request UserCreateRequest) (*User, error)
	GetUserByID(ctx context.Context, userId string) (*User, error)
	GetUserByAuthID(ctx context.Context, authId string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateAuthID(ctx context.Context, userId, authId string) error
	UpdateFirstName(ctx context.Context, userId, firstName string) error
	UpdateLastName(ctx context.Context, userId, lastName string) error
	UpdateEmail(ctx context.Context, userId, email string) error
	UpdateEmailVerificationStatus(ctx context.Context, email string, verified bool) error
	UpdatePhoneNumber(ctx context.Context, userId, phoneNumber string) error
	UpdatePhoneNumberVerificationStatus(ctx context.Context, phoneNumber string, verified bool) error
}

type repository struct {
	Pool sqlx.DB
}

func NewRepository(db sqlx.DB) Repository {
	return repository{
		Pool: db,
	}
}
