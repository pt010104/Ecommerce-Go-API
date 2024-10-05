package httpserver

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func (srv HTTPServer) Run() error {
	err := srv.mapHandlers()
	if err != nil {
		return err
	}

	ctx := context.Background()
	go func() {
		srv.gin.Run(fmt.Sprintf(":%d", srv.port))
	}()

	srv.l.Infof(ctx, "HTTP server is running on port %d", srv.port)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	srv.l.Infof(ctx, "Received signal: %v", <-ch)
	srv.l.Infof(ctx, "Shutting down HTTP server...")

	return nil
}
