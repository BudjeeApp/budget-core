package server

import (
	"github.com/BudjeeApp/budget-core/internal/category/repository"
	"github.com/BudjeeApp/budget-core/internal/config"
	"go.uber.org/zap"
)

var (
	InternalError = "internal server error"
)

type Server struct {
	Repo   repository.Repository
	Logger Logger
}

// TODO: Move this outside of category server
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

func NewServer(r repository.Repository) *Server {
	return &Server{
		r,
		config.Logger,
	}
}
