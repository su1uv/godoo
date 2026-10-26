package task

import (
	"context"

	db "github.com/su1uv/godoo/internal/db/sqlc"
)

type Service struct {
	repo *db.Queries
}

func NewService(repo *db.Queries) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTask(ctx context.Context, task *Task) (*Task, error) {

}
