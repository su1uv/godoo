package httpserver

import (
	"net/http"

	"github.com/su1uv/godoo/internal/helpers"
	"github.com/su1uv/godoo/internal/httpserver/middleware"
	"github.com/su1uv/godoo/internal/task"
)

func (a *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthz)

	taskService := task.NewService(a.queries)
	taskHandler := task.NewHandler(taskService, a.logger)

	mux.HandleFunc("POST /task", taskHandler.CreateTask)

	return middleware.Logger(a.logger)(mux)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	helpers.Encode(w, r, http.StatusOK, map[string]string{"status": "ok"})
}
