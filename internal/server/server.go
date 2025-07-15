package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	sloghttp "github.com/samber/slog-http"
)

type Server struct {
	*http.Server
}

func New(
	ctx context.Context,
	config Config,
	extraHandlers ...func() (string, http.Handler),
) *Server {
	mux := http.NewServeMux()

	s := &Server{
		Server: &http.Server{
			BaseContext: func(_ net.Listener) context.Context { return ctx },
			Addr:        fmt.Sprintf(":%d", config.Port),
			Handler:     mux,
		},
	}

	routes := map[string]http.Handler{
		"/": http.HandlerFunc(okHandler),
	}

	for _, extra := range extraHandlers {
		path, handler := extra()
		routes[path] = handler
	}

	middlewares := []func(http.Handler) http.Handler{
		sloghttp.New(slog.Default()),
		sloghttp.Recovery,
	}

	for path, handler := range routes {
		for _, middleware := range middlewares {
			handler = middleware(handler)
		}
		mux.Handle(path, handler)
	}

	mux.Handle("/probes/live", http.HandlerFunc(okHandler))
	mux.Handle("/probes/ready", http.HandlerFunc(okHandler))

	return s
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
