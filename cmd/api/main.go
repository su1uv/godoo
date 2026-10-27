package main

import (
	"os"

	"github.com/su1uv/godoo/internal/httpserver"
)

func main() {
	if err := httpserver.Run(); err != nil {
		os.Exit(1)
	}
}
