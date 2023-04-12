package server

import (
	"github.com/BudjeeApp/budget-core/config"
	"github.com/BudjeeApp/budget-core/internal/category/repository"
	"go.uber.org/zap"
)

var (
	InternalError = "internal server error"
)

type Server struct {
	Repo   repository.Repository
	Logger *zap.Logger
}

func NewServer(r repository.Repository) *Server {
	return &Server{
		r,
		config.Logger,
	}
}
