package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gaus57/http-multiplexer/server/handlers"
	"github.com/gaus57/http-multiplexer/server/middleware"
)

type Config struct {
	Addr          string // server address
	RequestsLimit int64  // the number of simultaneously executed requests
}

type Server struct {
	config *Config
	mp     handlers.Multiplexer
}

func New(cfg *Config, mp handlers.Multiplexer) *Server {
	return &Server{
		config: cfg,
		mp:     mp,
	}
}

// Serve start http server
func (s *Server) Serve(ctx context.Context) {
	sm := http.NewServeMux()
	sm.HandleFunc("/", middleware.LimitRequests(s.config.RequestsLimit, handlers.Home(s.mp)))

	httpServer := http.Server{
		Addr:    s.config.Addr,
		Handler: sm,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server listen error: %v\n", err)
		}
	}()

	<-ctx.Done()

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server shutdown failed: %v\n", err)
	}

	log.Println("server exited properly")
}
