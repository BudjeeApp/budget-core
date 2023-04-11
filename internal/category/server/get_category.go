package server

import (
	"context"
	"database/sql"
	"github.com/BudjeeApp/budget-core/internal/helpers"
	pb "github.com/BudjeeApp/budget-core/rpc/category"
	"github.com/twitchtv/twirp"
)

func (s *Server) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.Category, error) {
	if req.CategoryId == "" {
		return nil, twirp.RequiredArgumentError("category_id")
	}
	if !helpers.IsValidUUID(req.CategoryId) {
		return nil, twirp.InvalidArgumentError("category_id", "is an invalid uuid")
	}
	category, err := s.Repo.GetCategory(ctx, req.CategoryId)
	if err == sql.ErrNoRows {
		return nil, twirp.NotFoundError("category not found")
	} else if err != nil {
		return nil, twirp.InternalError(InternalError)
	}
	categoryResponse := category.ToProto()
	return &categoryResponse, nil
}
