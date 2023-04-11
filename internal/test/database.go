package test

import (
	"database/sql"
	"fmt"
	category "github.com/BudjeeApp/budget-core/internal/category/repository"
	user "github.com/BudjeeApp/budget-core/internal/user/repository"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"os"
	"sync"
	"testing"
)

type Repository struct {
	Category category.Repository
	User     user.Repository
}

var (
	oneTestDbPool sync.Once
	DB            *sqlx.DB
	Repo          *Repository
)

func GetDBInstance() (*Repository, error) {
	if DB == nil {
		var err error
		oneTestDbPool.Do(
			func() {
				var openErr error
				DB, openErr = sqlx.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
				if openErr != nil {
					err = fmt.Errorf("could not create database handle: %s", openErr.Error())
					return
				}
				if pingErr := DB.Ping(); pingErr != nil {
					err = fmt.Errorf("could not connect to database: %s", pingErr.Error())
					return
				}
				Repo = &Repository{
					Category: category.NewRepository(*DB),
					User:     user.NewRepository(*DB),
				}
			},
		)
		if err != nil {
			return nil, err
		}
	}
	return Repo, nil
}

func SeedUser(t *testing.T) *user.User {
	t.Helper()
	newUser := &user.User{}
	statement := `
		INSERT INTO users
			(auth_id, email, first_name, last_name, phone_number)
		VALUES 
			($1, $2, $3, $4, $5)
		RETURNING
			id, auth_id, email, email_verified, first_name, last_name, phone_number, phone_number_verified, created_at, updated_at
	`
	err := DB.QueryRowx(
		statement,
		uuid.NewString(), faker.Email(), faker.FirstName(), faker.LastName(), faker.Phonenumber(),
	).StructScan(newUser)
	if err != nil {
		t.Fatal("failed to seed user:" + err.Error())
	}
	return newUser
}

func RemoveUsers(t *testing.T) {
	t.Helper()
	if err := DB.Get(&struct{}{}, "DELETE FROM users"); err != nil && err != sql.ErrNoRows {
		t.Fatal("failed to remove users from database: " + err.Error())
	}
}

func SeedCategory(t *testing.T, userId string) *category.Category {
	t.Helper()
	newCategory := &category.Category{}
	statement := `
		INSERT INTO categories
			(user_id, name, allowance, cycle_type)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			id, user_id, name, parent_category_id, allowance, cycle_type, rollover, linked_users_id, created_at, updated_at
	`
	err := DB.QueryRowx(
		statement,
		userId, faker.Word(), 10000, "monthly",
	).StructScan(newCategory)
	if err != nil {
		t.Fatal("failed to seed category:" + err.Error())
	}
	return newCategory
}

func RemoveCategories(t *testing.T) {
	t.Helper()
	if err := DB.Get(&struct{}{}, "DELETE FROM categories"); err != nil && err != sql.ErrNoRows {
		t.Fatal("failed to remove categories from database: " + err.Error())
	}
}
