package task

import (
	"log/slog"
	"net/http"
	"time"
	"uuid"

	"github.com/su1uv/godoo/internal/db"
	"github.com/su1uv/godoo/internal/helpers"
)

type Handler struct {
	svc    *Service
	logger *slog.Logger
}

func NewHandler(svc *Service, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

type CreateTaskRequestDto struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueAt       time.Time `json:"due_at"`
}

type CreateTaskResponseDto struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueAt       *time.Time `json:"due_at"`
	CompletedAt *time.Time `json:"completed_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	payload, err := helpers.Decode[CreateTaskRequestDto](r)
	if err != nil {
		h.logger.Error("cannot read from body", "error", err)
		helpers.Encode(w, r, http.StatusBadRequest, err.Error())
		return
	}

	createdTask, err := h.svc.CreateTask(r.Context(), db.CreateTaskParams{
		Title:       payload.Title,
		Description: payload.Description,
		DueAt:       &payload.DueAt,
	})
	if err != nil {
		h.logger.Error("cannot create task", "error", err)
		helpers.Encode(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Encode(w, r, http.StatusCreated, CreateTaskResponseDto{
		ID:          createdTask.ID,
		Title:       createdTask.Title,
		Description: createdTask.Description,
		DueAt:       createdTask.DueAt,
		CompletedAt: createdTask.CompletedAt,
		DeletedAt:   createdTask.DeletedAt,
		CreatedAt:   createdTask.CreatedAt,
		UpdatedAt:   createdTask.UpdatedAt,
	})
}
