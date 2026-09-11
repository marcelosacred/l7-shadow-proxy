package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("health check: ok")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})


}