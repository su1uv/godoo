package task

import "net/http"

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type CreateTaskPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueAt       string `json:"due_at"`
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	task := &Task{}
	if err := h.svc.CreateTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
