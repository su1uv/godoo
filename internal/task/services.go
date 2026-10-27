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

func (s *Service) CreateTask(ctx context.Context, input *CreateTaskPayload) (*Task, error) {
	createdTask, err := s.repo.CreateTask(ctx, db.CreateTaskParams{
		Title:       input.Title,
		Description: input.Description,
		DueAt:       input.DueAt,
	})
	if err != nil {
		return nil, err
	}
	return &Task{
		ID:          createdTask.ID,
		CreatedAt:   createdTask.CreatedAt,
		UpdatedAt:   createdTask.UpdatedAt,
		Title:       createdTask.Title,
		Description: createdTask.Description,
		CompletedAt: &createdTask.CompletedAt,
		DueAt:       &createdTask.DueAt,
		DeletedAt:   &createdTask.DeletedAt,
	}, nil
}
