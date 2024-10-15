package httpserver

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/middleware"
)

func (srv HTTPServer) Run() error {
	ctx := context.Background()

	srv.gin = gin.Default()

	srv.gin.Use(middleware.CORSMiddleware())

	err := srv.mapHandlers()
	if err != nil {
		return err
	}

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
