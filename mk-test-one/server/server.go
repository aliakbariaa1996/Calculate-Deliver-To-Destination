package server

import (
	"context"
	"github.com/aliakbariaa1996/mk-test-one/config"
	"github.com/aliakbariaa1996/mk-test-one/internal/api/v1"
	loggerx "github.com/aliakbariaa1996/mk-test-one/internal/common/log"
	httpx "github.com/aliakbariaa1996/mk-test-one/internal/http"
	"log"
	"os"
	"os/signal"
	"time"
)

func RunServer(cfg *config.Config, logger *loggerx.Logger) error {
	// HTTP Server
	router := httpx.InitRouter()
	server, err := v1.NewServer(router, cfg, logger)
	if err != nil {
		return err
	}
	server.Server.Addr = ":" + cfg.Port
	server.Server.Handler = router
	server.Server.ReadTimeout = 10 * time.Second
	server.Server.WriteTimeout = 10 * time.Second
	server.Server.MaxHeaderBytes = 1 << 20
	go func() {
		if err := server.Server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return server.Server.Shutdown(ctx)
}