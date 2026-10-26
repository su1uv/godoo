package httpserver

import (
	"net/http"

	"github.com/su1uv/godoo/internal/platform"
	"github.com/su1uv/godoo/internal/task"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthz)

	taskService := task.NewService(taskRepo)
	taskHandler := task.NewHandler(taskService)

	mux.HandleFunc("POST /task", taskHandler.CreateTask)

	return mux
}

func healthz(w http.ResponseWriter, r *http.Request) {
	platform.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
