package main

import "net/http"

func registerRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	return mux
}
