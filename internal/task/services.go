package task

import (
	"context"

	"github.com/su1uv/godoo/internal/db"
)

type Service struct {
	repo *db.Queries
}

func NewService(repo *db.Queries) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTask(ctx context.Context, args db.CreateTaskParams) (db.Task, error) {
	createdTask, err := s.repo.CreateTask(ctx, args)
	if err != nil {
		return db.Task{}, err
	}

	return createdTask, nil
}
