package server

import (
	"context"
	"github.com/BudjeeApp/budget-core/internal/category/repository"
	"github.com/BudjeeApp/budget-core/internal/helpers"
	pb "github.com/BudjeeApp/budget-core/rpc/category"
	"github.com/google/uuid"
	"github.com/twitchtv/twirp"
)

func (s *Server) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error) {
	if req.Name == "" {
		return nil, twirp.RequiredArgumentError("name")
	}
	if req.UserId == "" {
		return nil, twirp.RequiredArgumentError("user_id")
	}
	if _, err := uuid.Parse(req.UserId); err != nil {
		return nil, twirp.InvalidArgumentError("user_id", "is an invalid uuid")
	}
	parentCategoryId, err := helpers.NullStringFromUUID("parent_category_id", req.GetParentCategoryId())
	if err != nil {
		return nil, err
	}
	linkedUsersId, err := helpers.NullStringFromUUID("linked_users_id", req.GetLinkedUsersId())
	if err != nil {
		return nil, err
	}
	repoCategoryCreateRequest := repository.CategoryCreateRequest{
		UserID:           req.GetUserId(),
		Name:             req.GetName(),
		ParentCategoryID: parentCategoryId,
		Allowance:        req.GetAllowance(),
		CycleType:        req.GetCycleType().String(),
		Rollover:         req.GetRollover(),
		LinkedUsersID:    linkedUsersId,
	}
	createdCategory, err := s.Repo.CreateCategory(ctx, repoCategoryCreateRequest)
	if err != nil {
		return nil, twirp.InternalError(InternalError)
	}
	categoryResponse := createdCategory.ToProto()
	return &categoryResponse, nil
}
