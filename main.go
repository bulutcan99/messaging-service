package main

import (
	"context"
	"os/signal"
	"websocket-azure/presentation/http"
	"websocket-azure/shared/logger"
	"websocket-azure/shared/setting"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

func main() {
	if _, err := maxprocs.Set(); err != nil {
		panic("failed to set max procs: " + err.Error())
	}

	ctx, cancel := signal.NotifyContext(context.Background(), setting.InterruptSignals...)
	defer cancel()

	log := logger.CreateLogger(1)
	defer log.Sync()
	zap.ReplaceGlobals(log)

	server := http.NewServer(ctx)
	if err := server.Start(ctx); err != nil {
		panic("Server starting error!")
	}

	zap.S().Debug("⚡ Service name: ", "WEBSOCKET CHAT")

	<-ctx.Done()
	zap.S().Debug("Context signal received, shutting down")
	zap.S().Debug("Waiting for all goroutines to finish")
	zap.S().Debug("Shutting down successfully")
}
