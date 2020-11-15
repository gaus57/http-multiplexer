package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gaus57/http-multiplexer/server"
	"github.com/gaus57/http-multiplexer/services/client"
	"github.com/gaus57/http-multiplexer/services/multiplexer"
)

// defaultConfig default app configuration
var defaultConfig = struct {
	server *server.Config
	mp     *multiplexer.Config
}{
	server: &server.Config{
		Addr:          ":8080",
		RequestsLimit: 100,
	},
	mp: &multiplexer.Config{
		RequestsLimit:  4,
		RequestTimeOut: time.Second,
	},
}

func main() {
	port := flag.Int("p", 8080, "Port for http server")
	flag.Parse()
	if port != nil {
		defaultConfig.server.Addr = fmt.Sprintf(":%d", *port)
	}

	ctx := graceContext()

	mp := multiplexer.New(
		defaultConfig.mp,
		client.New(&http.Client{}),
	)

	srv := server.New(
		defaultConfig.server,
		mp,
	)

	srv.Serve(ctx)
}

func graceContext() context.Context {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-done
		cancel()
	}()

	return ctx
}
