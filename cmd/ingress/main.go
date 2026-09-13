package main

import (
	"errors"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
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

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).DialContext,
		ResponseHeaderTimeout: 3 * time.Second,
		MaxConnsPerHost: 256,
		MaxIdleConnsPerHost: 64,
		IdleConnTimeout: 90 * time.Second,
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func (pr *httputil.ProxyRequest)  {
			pr.SetURL(target)
		},
		Transport: transport,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			status := http.StatusBadGateway // 502 (по умолчанию)

			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				status = http.StatusGatewayTimeout // 504
			}

			logger.Error("upstream unreachable",
				"path", r.URL.Path,
				"method", r.Method,
				"status", status,
				"err", err,
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
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