package repository_test

import (
	"context"
	"database/sql"
	"github.com/BudjeeApp/budget-core/internal/test"
	. "github.com/BudjeeApp/budget-core/internal/user/repository"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_CreateUser(t *testing.T) {
	t.Run("should create a user with minimum fields", func(t *testing.T) {
		setupTest(t)
		defer teardownTest(t)
		req := UserCreateRequest{
			AuthID:    uuid.NewString(),
			Email:     faker.Email(),
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
		}
		user, err := repo.CreateUser(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, req.AuthID, user.AuthID)
		assert.Equal(t, req.Email, user.Email)
		assert.Equal(t, req.FirstName, user.FirstName)
		assert.Equal(t, req.LastName, user.LastName)
		assert.Equal(t, false, user.PhoneNumber.Valid, "expected new user to have null phone_number")
	})

	t.Run("should create a user with all fields", func(t *testing.T) {
		setupTest(t)
		defer teardownTest(t)
		req := UserCreateRequest{
			AuthID:      uuid.NewString(),
			Email:       faker.Email(),
			FirstName:   faker.FirstName(),
			LastName:    faker.LastName(),
			PhoneNumber: sql.NullString{String: faker.Phonenumber(), Valid: true},
		}
		user, err := repo.CreateUser(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, req.AuthID, user.AuthID)
		assert.Equal(t, req.Email, user.Email)
		assert.Equal(t, req.FirstName, user.FirstName)
		assert.Equal(t, req.LastName, user.LastName)
		assert.Equal(t, req.PhoneNumber.String, user.PhoneNumber.String)
	})

	t.Run("should create a user with unverified email", func(t *testing.T) {
		setupTest(t)
		defer teardownTest(t)
		req := UserCreateRequest{
			AuthID:    uuid.NewString(),
			Email:     faker.Email(),
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
		}
		user, err := repo.CreateUser(context.Background(), req)
		assert.NoError(t, err)
		assert.False(t, user.EmailVerified, "expected new user to have an unverified email")
	})

	t.Run("should create a user with unverified phone number", func(t *testing.T) {
		setupTest(t)
		defer teardownTest(t)
		req := UserCreateRequest{
			AuthID:      uuid.NewString(),
			Email:       faker.Email(),
			FirstName:   faker.FirstName(),
			LastName:    faker.LastName(),
			PhoneNumber: sql.NullString{String: faker.Phonenumber(), Valid: true},
		}
		user, err := repo.CreateUser(context.Background(), req)
		assert.NoError(t, err)
		assert.False(t, user.PhoneNumberVerified, "expected new user to have an unverified phone number")
	})
}

func TestRepository_GetUserByID(t *testing.T) {
	withUserSeed(t, func(t *testing.T, user *User) {
		t.Run("should get a user if user ID exists", func(t *testing.T) {
			gotUser, err := repo.GetUserByID(context.Background(), user.ID)
			assert.NoError(t, err)
			assert.Equal(t, user.ID, gotUser.ID)
		})

		t.Run("should return an error if user ID is not found", func(t *testing.T) {
			gotUser, err := repo.GetUserByID(context.Background(), uuid.NewString())
			assert.Error(t, err)
			assert.Nil(t, gotUser)
		})
	})

}

func TestRepository_GetUserByAuthID(t *testing.T) {
	withUserSeed(t, func(t *testing.T, user *User) {
		t.Run("should get a user if auth ID exists", func(t *testing.T) {
			gotUser, err := repo.GetUserByAuthID(context.Background(), user.AuthID)
			assert.NoError(t, err)
			assert.Equal(t, user.ID, gotUser.ID)
		})

		t.Run("should return an error if user ID is not found", func(t *testing.T) {
			gotUser, err := repo.GetUserByAuthID(context.Background(), uuid.NewString())
			assert.Error(t, err)
			assert.Nil(t, gotUser)
		})
	})

}

func TestRepository_GetUserByEmail(t *testing.T) {
	withUserSeed(t, func(t *testing.T, user *User) {
		t.Run("should get a user if email exists", func(t *testing.T) {
			gotUser, err := repo.GetUserByEmail(context.Background(), user.Email)
			assert.NoError(t, err)
			assert.Equal(t, user.ID, gotUser.ID)
		})

		t.Run("should return an error if email is not found", func(t *testing.T) {
			gotUser, err := repo.GetUserByEmail(context.Background(), "nonexistant@email.com")
			assert.Error(t, err)
			assert.Nil(t, gotUser)
		})
	})
}

func TestRepository_UpdateAuthID(t *testing.T) {
	withUserSeed(t, func(t *testing.T, user *User) {
		t.Run("should be successful if user ID exists and auth ID is valid", func(t *testing.T) {
			newUUID := uuid.NewString()
			err := repo.UpdateAuthID(context.Background(), user.ID, newUUID)
			assert.NoError(t, err)
			fetchUser(t, user)
			assert.Equal(t, newUUID, user.AuthID)
		})

		t.Run("should return an error if user ID is not found", func(t *testing.T) {
			newUUID := uuid.NewString()
			err := repo.UpdateAuthID(context.Background(), uuid.NewString(), newUUID)
			assert.Error(t, err)
		})

		t.Run("should return an error if auth ID is not a valid UUID", func(t *testing.T) {
			newUUID := "invalid-UUID"
			err := repo.UpdateAuthID(context.Background(), user.ID, newUUID)
			assert.Error(t, err)
		})

		t.Run("should return an error if given auth ID is taken", func(t *testing.T) {
			anotherUser := test.SeedUser(t)
			err := repo.UpdateAuthID(context.Background(), user.ID, anotherUser.AuthID)
			assert.Error(t, err)
		})
	})
}

func TestRepository_UpdateFirstName(t *testing.T) {
	withUserSeed(t, func(t *testing.T, user *User) {
		t.Run("should be successful if user ID exists", func(t *testing.T) {
			err := repo.UpdateFirstName(context.Background(), user.ID, "Fred")
			assert.NoError(t, err)
			fetchUser(t, user)
			assert.Equal(t, "Fred", user.FirstName)
		})

		t.Run("should return an error if user ID is not found", func(t *testing.T) {
			err := repo.UpdateFirstName(context.Background(), uuid.NewString(), "Fred")
			assert.Error(t, err)
		})
	})
}

func TestRepository_UpdateLastName(t *testing.T) {
	withUserSeed(t, func(t *testing.T, user *User) {
		t.Run("should be successful if user ID exists", func(t *testing.T) {
			err := repo.UpdateLastName(context.Background(), user.ID, "Jones")
			assert.NoError(t, err)
			fetchUser(t, user)
			assert.Equal(t, "Jones", user.LastName)
		})

		t.Run("should return an error if user ID is not found", func(t *testing.T) {
			err := repo.UpdateLastName(context.Background(), uuid.NewString(), "Jones")
			assert.Error(t, err)
		})
	})
}

func TestRepository_UpdateEmail(t *testing.T) {
	withUserSeed(t, func(t *testing.T, user *User) {
		t.Run("should be successful if user ID exists", func(t *testing.T) {
			err := repo.UpdateEmail(context.Background(), user.ID, "test@email.com")
			assert.NoError(t, err)
			fetchUser(t, user)
			assert.Equal(t, "test@email.com", user.Email)
		})

		t.Run("should return an error if user ID is not found", func(t *testing.T) {
			err := repo.UpdateEmail(context.Background(), uuid.NewString(), "test2@email.com")
			assert.Error(t, err)
		})

		t.Run("should return an error if given email is already taken", func(t *testing.T) {
			anotherUser := test.SeedUser(t)
			err := repo.UpdateEmail(context.Background(), user.ID, anotherUser.Email)
			assert.Error(t, err)
		})
	})
}

func TestRepository_UpdateEmailVerificationStatus(t *testing.T) {
	withUserSeed(t, func(t *testing.T, user *User) {
		t.Run("should be successful if email exists", func(t *testing.T) {
			err := repo.UpdateEmailVerificationStatus(context.Background(), user.Email, true)
			assert.NoError(t, err)
			fetchUser(t, user)
			assert.Equal(t, true, user.EmailVerified, "expected email_verified to be updated to true")
		})

		t.Run("should return an error if email is not found", func(t *testing.T) {
			err := repo.UpdateEmailVerificationStatus(context.Background(), "nonexistant@email.com", true)
			assert.Error(t, err)
		})
	})
}

func TestRepository_UpdatePhoneNumber(t *testing.T) {
	withUserSeed(t, func(t *testing.T, user *User) {
		t.Run("should be successful if user ID exists", func(t *testing.T) {
			err := repo.UpdatePhoneNumber(context.Background(), user.ID, "1231231233")
			assert.NoError(t, err)
			fetchUser(t, user)
			assert.Equal(t, "1231231233", user.PhoneNumber.String)
		})

		t.Run("should return an error if user ID is not found", func(t *testing.T) {
			err := repo.UpdatePhoneNumber(context.Background(), uuid.NewString(), "00000000")
			assert.Error(t, err)
		})

		t.Run("should return an error if given phone number is already taken", func(t *testing.T) {
			anotherUser := test.SeedUser(t)
			err := repo.UpdatePhoneNumber(context.Background(), anotherUser.ID, "9999999999")
			assert.NoError(t, err)
			err = repo.UpdatePhoneNumber(context.Background(), user.ID, "9999999999")
			assert.Error(t, err)
		})
	})
}

func TestRepository_UpdatePhoneNumberVerificationStatus(t *testing.T) {
	withUserSeed(t, func(t *testing.T, user *User) {
		err := repo.UpdatePhoneNumber(context.Background(), user.ID, "1231231233")
		if err != nil {
			t.Fatal("failed updating user phone number: " + err.Error())
		}

		t.Run("should be successful if phone number exists", func(t *testing.T) {
			err = repo.UpdatePhoneNumberVerificationStatus(context.Background(), "1231231233", true)
			assert.NoError(t, err)
			fetchUser(t, user)
			assert.Equal(t, true, user.PhoneNumberVerified, "expected phone_number_verified to be updated to true")
		})

		t.Run("should return an error if phone number is not found", func(t *testing.T) {
			err = repo.UpdatePhoneNumberVerificationStatus(context.Background(), "0000000000", true)
			assert.Error(t, err)
		})
	})
}

func withUserSeed(t *testing.T, testFunc func(t *testing.T, user *User)) {
	setupTest(t)
	defer teardownTest(t)
	user := test.SeedUser(t)
	testFunc(t, user)
}

func fetchUser(t *testing.T, user *User) {
	statement := `
			SELECT
				*
			FROM
				users
			WHERE id = $1
		`
	err := test.DB.Get(user, statement, user.ID)
	if err != nil {
		t.Fatal(err.Error())
	}
}
