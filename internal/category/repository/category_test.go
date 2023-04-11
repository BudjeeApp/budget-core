package repository_test

import (
	"context"
	"database/sql"
	. "github.com/BudjeeApp/budget-core/internal/category/repository"
	"github.com/BudjeeApp/budget-core/internal/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_CreateCategory(t *testing.T) {
	assert := assert.New(t)

	t.Run("should create a category with minimum fields", func(t *testing.T) {
		setupTest(t)
		defer teardownTest(t)
		user := test.SeedUser(t)
		req := CategoryCreateRequest{
			UserID:    user.ID,
			Name:      "testName",
			CycleType: "monthly",
		}
		category, err := repo.CreateCategory(context.Background(), req)
		assert.NoError(err)
		assert.Equal(req.Name, category.Name)
		assert.Equal(req.UserID, category.UserID)
	})

	t.Run("should create a category with all fields", func(t *testing.T) {
		t.Skip("Need to add linked users stuff")
		setupTest(t)
		defer teardownTest(t)
		parentCategoryRequest := CategoryCreateRequest{
			UserID:    uuid.NewString(),
			Name:      "testName",
			CycleType: "monthly",
		}
		parentCategory, err := repo.CreateCategory(context.Background(), parentCategoryRequest)
		assert.NoError(err)

		subCategoryRequest := CategoryCreateRequest{
			UserID:           uuid.NewString(),
			Name:             "Sub Category",
			ParentCategoryID: parentCategory.ParentCategoryID,
			Allowance:        10032,
			CycleType:        "weekly",
			Rollover:         true,
			LinkedUsersID:    sql.NullString{String: uuid.NewString(), Valid: true},
		}
		subCategory, err := repo.CreateCategory(context.Background(), subCategoryRequest)
		assert.NoError(err)
		assert.Equal(subCategoryRequest.Name, subCategory.Name)
		assert.Equal(subCategoryRequest.UserID, subCategory.UserID)
		assert.Equal(subCategoryRequest.ParentCategoryID, subCategory.ParentCategoryID)
		assert.Equal(subCategoryRequest.Allowance, subCategory.Allowance)
		assert.Equal(subCategoryRequest.CycleType, subCategory.CycleType)
		assert.Equal(subCategoryRequest.Rollover, subCategory.Rollover)
		assert.Equal(subCategoryRequest.LinkedUsersID, subCategory.LinkedUsersID)
	})

	t.Run("should require an existing user ID", func(t *testing.T) {
		setupTest(t)
		defer teardownTest(t)
		req := CategoryCreateRequest{
			UserID:    uuid.NewString(),
			Name:      "testName",
			CycleType: "monthly",
		}
		category, err := repo.CreateCategory(context.Background(), req)
		assert.Error(err)
		assert.Nil(category)
	})

	t.Run("should require an existing CategoryID for the ParentCategoryID", func(t *testing.T) {
		setupTest(t)
		defer teardownTest(t)
		user := test.SeedUser(t)
		req := CategoryCreateRequest{
			UserID:           user.ID,
			Name:             "Nonexistent Parent Category",
			ParentCategoryID: sql.NullString{String: uuid.NewString(), Valid: true},
		}
		category, err := repo.CreateCategory(context.Background(), req)
		assert.Error(err)
		assert.Nil(category)
	})
}

func TestRepository_GetCategory(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)
	user := test.SeedUser(t)

	t.Run("should retrieve the correct category", func(t *testing.T) {
		category := test.SeedCategory(t, user.ID)

		retrievedCategory, err := repo.GetCategory(context.Background(), category.ID)
		assert.NoError(t, err)
		assert.Equal(t, category.ID, retrievedCategory.ID)
	})

	t.Run("should return an error if category is not found", func(t *testing.T) {
		retrievedCategory, err := repo.GetCategory(context.Background(), uuid.NewString())
		assert.Nil(t, retrievedCategory)
		assert.Error(t, err)
	})
}
