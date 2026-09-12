package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type response struct {
	Name    string      `json:"name"`
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Headers http.Header `json:"headers"`
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	listen 	  := flag.String("listen", ":8080", "adress to listen on")
	name 	  := flag.String("name",  "service", "service name")
	delay 	  := flag.Duration("delay", 0, "response delay")
	jitter 	  := flag.Duration("jitter", 0, "random delay jitter")
	errorRate := flag.Float64( "errorRate", 0, "error rate")

	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("healthz check: ok")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sleep := *delay
		if *jitter > 0 { sleep += time.Duration(rand.Int63n(int64(*jitter))) }
		time.Sleep(sleep)

		if rand.Float64() < *errorRate {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("error"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response{
			Name:    *name,
			Method:  r.Method,
			Path:    r.URL.Path,
			Headers: r.Header,
		})

		logger.Info("request processed",
			"service", *name,
			"delay", sleep,
			"status", http.StatusOK,
		)	})

	err := http.ListenAndServe(*listen, mux)
	if err != nil {
		logger.Error(err.Error())
	}

}