package main

import (
	"flag"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("l7-shadow-proxy starting...")

	// ---------------------------------
	listen := flag.String("listen", ":8080", "adress to listen on")


	flag.Parse()
	// ---------------------------------

	const backendURL = "http://localhost:8081"

	target, err := url.Parse(backendURL)
	if err != nil {
		logger.Error("parsing backend url", "err", err)
		os.Exit(1)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func (pr *httputil.ProxyRequest)  {
			pr.SetURL(target)
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Error("upstream unreachable", "path", r.URL.Path, "method", r.Method, "err", err)
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(`{"error":"upstream unreachable"}`))
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/", proxy)

	err = http.ListenAndServe(*listen, mux)
	if err != nil {
		logger.Error("starting: ", "err", err)
		os.Exit(1)
	}
}