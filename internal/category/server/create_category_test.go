package server

import (
	"context"
	"database/sql"
	"github.com/BudjeeApp/budget-core/internal/category/mocks"
	"github.com/BudjeeApp/budget-core/internal/category/repository"
	pb "github.com/BudjeeApp/budget-core/rpc/category"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_CreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(ctrl)
	defer ctrl.Finish()

	server := NewServer(mockRepo)
	twirpHandler := pb.NewCategoryServiceServer(server)
	testServer := httptest.NewServer(twirpHandler)
	defer testServer.Close()

	client := pb.NewCategoryServiceProtobufClient(testServer.URL, http.DefaultClient)

	t.Run("should create the category for the user who requested it", func(t *testing.T) {
		t.Skip("Pending. May need to move this test")
	})

	t.Run("should respond with the generated category uuid", func(t *testing.T) {
		mockRepo.
			EXPECT().
			CreateCategory(gomock.Any(), gomock.Any()).
			Return(successfulCategory, nil).
			Times(1)
		req := &pb.CreateCategoryRequest{
			UserId:           uuid.NewString(),
			Name:             "Some Test Name",
			ParentCategoryId: "dd684402-9638-4576-9fdb-823688f44ff9",
			Allowance:        60000,
		}
		resp, err := client.CreateCategory(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, successfulCategory.ID, resp.Id)
	})

	t.Run("should respond with the correct cycle type if cycle type is given", func(t *testing.T) {
		var (
			expectedCategory = &repository.Category{
				ID:        uuid.NewString(),
				UserID:    uuid.NewString(),
				Name:      "Some Test Name",
				Allowance: 60000,
				CycleType: "quarterly",
				Rollover:  false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		)
		mockRepo.
			EXPECT().
			CreateCategory(gomock.Any(), gomock.Any()).
			Return(expectedCategory, nil).
			Times(1)
		req := &pb.CreateCategoryRequest{
			UserId:    uuid.NewString(),
			Name:      "Some Test Name",
			CycleType: pb.CycleType_quarterly,
		}
		resp, err := client.CreateCategory(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, "quarterly", resp.CycleType.String())
	})

	t.Run("should be able to accept only the category name and user ID in the request", func(t *testing.T) {
		mockRepo.
			EXPECT().
			CreateCategory(gomock.Any(), gomock.Any()).
			Return(successfulCategory, nil).
			Times(1)
		req := &pb.CreateCategoryRequest{
			UserId: uuid.NewString(),
			Name:   "Successful Category",
		}
		_, err := client.CreateCategory(context.Background(), req)
		assert.NoError(t, err)
	})

	repositoryErrorTestCases := []struct {
		TestName string
		RepoFunc func(repo *mocks.MockRepository)
		Expected string
	}{
		{
			TestName: "should return an internal error when the repository fails while creating the category",
			RepoFunc: func(repo *mocks.MockRepository) {
				repo.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Return(nil, sql.ErrConnDone).Times(1)
			},
			Expected: "internal",
		},
	}
	for _, tc := range repositoryErrorTestCases {
		t.Run(tc.TestName, func(t *testing.T) {
			tc.RepoFunc(mockRepo)
			req := &pb.CreateCategoryRequest{
				UserId: uuid.NewString(),
				Name:   "Some Category",
			}
			resp, err := client.CreateCategory(context.Background(), req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.Expected)
			assert.Nil(t, resp)

		})
	}

	invalidArgumentTestCases := []struct {
		TestName              string
		FieldName             string
		CreateCategoryRequest *pb.CreateCategoryRequest
	}{
		{
			TestName:  "empty Name",
			FieldName: "name",
			CreateCategoryRequest: &pb.CreateCategoryRequest{
				UserId:    uuid.NewString(),
				Allowance: 55000,
			},
		},
		{
			TestName:  "empty userId",
			FieldName: "user_id",
			CreateCategoryRequest: &pb.CreateCategoryRequest{
				Name:      "Category Name",
				Allowance: 55000,
			},
		},
		{
			TestName:  "invalid userId",
			FieldName: "user_id",
			CreateCategoryRequest: &pb.CreateCategoryRequest{
				UserId:    "123",
				Name:      "Category Name",
				Allowance: 55000,
			},
		},
		{
			TestName:  "invalid parentCategoryId",
			FieldName: "parent_category_id",
			CreateCategoryRequest: &pb.CreateCategoryRequest{
				UserId:           uuid.NewString(),
				Name:             "Category",
				ParentCategoryId: "123",
			},
		},
		{
			TestName:  "invalid linkedUsersId",
			FieldName: "linked_users_id",
			CreateCategoryRequest: &pb.CreateCategoryRequest{
				UserId:        uuid.NewString(),
				Name:          "Category",
				LinkedUsersId: "13a6682f-795c-49c1-bfbb-f94f40eef",
			},
		},
	}

	for _, tc := range invalidArgumentTestCases {
		t.Run("should return an invalid argument error for "+tc.TestName, func(t *testing.T) {
			mockRepo.
				EXPECT().
				CreateCategory(gomock.Any(), gomock.Any()).
				Times(0)
			req := tc.CreateCategoryRequest
			resp, err := client.CreateCategory(context.Background(), req)
			assert.Error(t, err, "expected an error")
			assert.Contains(t, err.Error(), "invalid_argument")
			assert.Contains(t, err.Error(), tc.FieldName)
			assert.Nil(t, resp)
		})
	}
}
