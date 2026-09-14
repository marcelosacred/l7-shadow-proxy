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

	"github.com/marcelosacred/l7-shadow-proxy/internal/config"
)

// TODO:
// Струтурировать все ошибки в одном файле для удобства поиска и исправления.

func buildProxy(up config.UpstreamConfig, logger *slog.Logger) (*httputil.ReverseProxy, error) {
	if len(up.Targets) == 0 { return nil, errors.New("upstream has no targets") }

	target, err := url.Parse(up.Targets[0]);
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: up.Timeout,
		}).DialContext,
		ResponseHeaderTimeout: 3 * time.Second,
		MaxConnsPerHost: up.MaxConnsPerHost,
		MaxIdleConnsPerHost: 64,
		IdleConnTimeout: 90 * time.Second,
	}

	return &httputil.ReverseProxy{
		Rewrite: func (pr *httputil.ProxyRequest)  {
			pr.SetURL(target)
		},
		Transport: transport,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
					status := http.StatusBadGateway
					var netErr net.Error
					if errors.As(err, &netErr) && netErr.Timeout() {
						status = http.StatusGatewayTimeout
					}
					logger.Error("upstream error",
						"path", r.URL.Path, "method", r.Method, "status", status, "err", err)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(status)
					_, _ = w.Write([]byte(`{"error":"upstream error"}`))
				},
	}, nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("l7-shadow-proxy starting...")

	// ---------------------------------
	listen := flag.String("listen", ":8080", "adress to listen on")
	configPath:= flag.String("config", "configs/example.yaml", "path to config")

	flag.Parse()
	// ---------------------------------

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("load config", "err", err)
		os.Exit(1)
	}

	proxies := make(map[string]*httputil.ReverseProxy)
	for name, up := range cfg.Upstreams {
		p, err := buildProxy(up, logger)
		if err != nil {
			logger.Error("build proxy", "upstream", name, "err", err)
			os.Exit(1)
		}
		proxies[name] = p
	}
	logger.Info("proxies built", "count", len(proxies))

	mux := http.NewServeMux()
	mux.Handle("/", proxies["backend1"])

	err = http.ListenAndServe(*listen, mux)
	if err != nil {
		logger.Error("starting: ", "err", err)
		os.Exit(1)
	}
}