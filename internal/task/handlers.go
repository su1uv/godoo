package task

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/su1uv/godoo/internal/helpers"
)

type Handler struct {
	svc    *Service
	logger *slog.Logger
}

func NewHandler(svc *Service, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

type CreateTaskPayload struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueAt       time.Time `json:"due_at"`
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var payload CreateTaskPayload
	if err := helpers.ReadJSON(w, r, &payload); err != nil {
		h.logger.Error("cannot read from body", "error", err)
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdTask, err := h.svc.CreateTask(r.Context(), &payload)
	if err != nil {
		h.logger.Error("cannot create task", "error", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, *createdTask)
}
