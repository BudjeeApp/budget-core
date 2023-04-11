package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID                  string         `db:"id"`
	AuthID              string         `db:"auth_id"`
	Email               string         `db:"email"`
	FirstName           string         `db:"first_name"`
	LastName            string         `db:"last_name"`
	PhoneNumber         sql.NullString `db:"phone_number"`
	EmailVerified       bool           `db:"email_verified"`
	PhoneNumberVerified bool           `db:"phone_number_verified"`
	CreatedAt           time.Time      `db:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at"`
}

type UserCreateRequest struct {
	AuthID      string         `db:"auth_id"`
	Email       string         `db:"email"`
	FirstName   string         `db:"first_name"`
	LastName    string         `db:"last_name"`
	PhoneNumber sql.NullString `db:"phone_number"`
}

func (repo repository) CreateUser(ctx context.Context, user UserCreateRequest) (*User, error) {
	if user.PhoneNumber.String == "" {
		user.PhoneNumber.Valid = false
	}
	result := &User{}
	statement := `
		INSERT INTO users
			(auth_id, email, first_name, last_name, phone_number)
		VALUES 
			($1, $2, $3, $4, $5)
		RETURNING
			id, auth_id, email, email_verified, first_name, last_name, phone_number, phone_number_verified, created_at, updated_at
	`
	err := repo.Pool.QueryRowxContext(
		ctx,
		statement,
		user.AuthID, user.Email, user.FirstName, user.LastName, user.PhoneNumber,
	).StructScan(result)
	if err != nil {
		return nil, fmt.Errorf("failed to create user %s: %w", user.AuthID, err)
	}
	return result, nil
}

func (repo repository) GetUserByID(ctx context.Context, userId string) (*User, error) {
	result := &User{}
	err := repo.getUserByColumn(ctx, result, "id", userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo repository) GetUserByAuthID(ctx context.Context, authId string) (*User, error) {
	result := &User{}
	err := repo.getUserByColumn(ctx, result, "auth_id", authId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	result := &User{}
	err := repo.getUserByColumn(ctx, result, "email", email)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo repository) UpdateAuthID(ctx context.Context, userId, authId string) error {
	return repo.updateUserByColumn(ctx, "id", userId, "auth_id", authId)
}

func (repo repository) UpdateFirstName(ctx context.Context, userId, firstName string) error {
	return repo.updateUserByColumn(ctx, "id", userId, "first_name", firstName)
}

func (repo repository) UpdateLastName(ctx context.Context, userId, lastName string) error {
	return repo.updateUserByColumn(ctx, "id", userId, "last_name", lastName)
}

func (repo repository) UpdateEmail(ctx context.Context, userId, email string) error {
	return repo.updateUserByColumn(ctx, "id", userId, "email", email)
}

func (repo repository) UpdateEmailVerificationStatus(ctx context.Context, email string, verified bool) error {
	return repo.updateUserByColumn(ctx, "email", email, "email_verified", verified)
}

func (repo repository) UpdatePhoneNumber(ctx context.Context, userId, phoneNumber string) error {
	return repo.updateUserByColumn(ctx, "id", userId, "phone_number", phoneNumber)
}

func (repo repository) UpdatePhoneNumberVerificationStatus(ctx context.Context, phoneNumber string, verified bool) error {
	return repo.updateUserByColumn(ctx, "phone_number", phoneNumber, "phone_number_verified", verified)
}

// getUserByColumn fetches the user by a given field and its value specified by 'columnName' and 'value' respectively.
// The retrieved user is stored in 'dest'.
// Returns an error if the operation fails.
func (repo repository) getUserByColumn(ctx context.Context, dest *User, columnName, value string) error {
	statement := `
		SELECT
			id, auth_id, email, first_name, last_name, phone_number, email_verified, phone_number_verified, created_at, updated_at
		FROM
			users
		WHERE
			%s = $1;
	`
	return repo.Pool.GetContext(ctx, dest, fmt.Sprintf(statement, columnName), value)
}

// updateUserByColumn updates a user's information by specifying a column to search by and a column to update.
// The function takes in the search column (by Column) with its value (byValue) and the update column (updateColumn)
// with its update Value (updateValue). It returns an error if the update operation fails or if the user was not found.
func (repo repository) updateUserByColumn(ctx context.Context, byColumn string, byValue interface{}, updateColumn string, updateValue interface{}) error {
	statement := `
		UPDATE
			users
		SET
			%s = $1
		WHERE
		    %s = $2
		RETURNING id;
	`
	result, err := repo.Pool.ExecContext(ctx, fmt.Sprintf(statement, updateColumn, byColumn), updateValue, byValue)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("no user found by column %s, value %v", byColumn, byValue)
	}
	return nil
}
