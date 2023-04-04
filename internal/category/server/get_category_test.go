package server

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jee-lee/budget-core/internal/category/mocks"
	"github.com/jee-lee/budget-core/internal/category/repository"
	pb "github.com/jee-lee/budget-core/rpc/category"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_GetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(ctrl)
	defer ctrl.Finish()

	server := NewServer(mockRepo)
	twirpHandler := pb.NewCategoryServiceServer(server)
	testServer := httptest.NewServer(twirpHandler)
	defer testServer.Close()

	client := pb.NewCategoryServiceProtobufClient(testServer.URL, http.DefaultClient)

	t.Run("should get the correct category", func(t *testing.T) {
		existingCategory := &repository.Category{
			ID:          uuid.MustParse("13a6682f-795c-49c1-bfbb-f94f4b770eef"),
			UserID:      uuid.MustParse("2b807819-078c-4d0d-b2b3-6204ff95f967"),
			Name:        "Successful Category",
			CycleTypeID: 2,
		}
		mockRepo.
			EXPECT().
			GetCategory(gomock.Any(), gomock.Any()).
			Do(func(context interface{}, id *uuid.UUID) {
				if _, ok := interface{}(id).(*uuid.UUID); !ok {
					t.Errorf("expected id to be of type *uuid.UUID, but got %T", id)
				}
				assert.Equal(t, id.String(), existingCategory.ID.String())
			}).
			Return(existingCategory, nil).
			Times(1)
		mockRepo.
			EXPECT().
			GetCycleTypeByID(gomock.Any(), 2).
			Return(defaultCycleType, nil).
			Times(1)
		req := &pb.GetCategoryRequest{
			CategoryId: "13a6682f-795c-49c1-bfbb-f94f4b770eef",
		}
		resp, err := client.GetCategory(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, existingCategory.ID.String(), resp.Id)
	})

	invalidArgumentTestCases := []struct {
		TestName              string
		CreateCategoryRequest *pb.GetCategoryRequest
	}{
		{
			TestName:              "empty categoryId",
			CreateCategoryRequest: &pb.GetCategoryRequest{},
		},
		{
			TestName: "invalid categoryId",
			CreateCategoryRequest: &pb.GetCategoryRequest{
				CategoryId: "123",
			},
		},
	}

	for _, tc := range invalidArgumentTestCases {
		t.Run("should return an invalid argument error for "+tc.TestName, func(t *testing.T) {
			mockRepo.
				EXPECT().
				GetCategory(gomock.Any(), gomock.Any).
				Times(0)
			resp, err := client.GetCategory(context.Background(), tc.CreateCategoryRequest)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid_argument")
			assert.Nil(t, resp)
		})
	}

	t.Run("should return an internal server error if the repository return an error getting the category", func(t *testing.T) {
		mockRepo.
			EXPECT().
			GetCategory(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, sql.ErrConnDone)
		req := &pb.GetCategoryRequest{
			CategoryId: "13a6682f-795c-49c1-bfbb-f94f4b770eef",
		}
		resp, err := client.GetCategory(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "internal")
		assert.Nil(t, resp)
	})

	t.Run("should return a not found error when the category does not exist", func(t *testing.T) {
		mockRepo.
			EXPECT().
			GetCategory(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, sql.ErrNoRows)
		req := &pb.GetCategoryRequest{
			CategoryId: "13a6682f-795c-49c1-bfbb-f94f4b770eef",
		}
		resp, err := client.GetCategory(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not_found")
		assert.Nil(t, resp)
	})
}
