package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infof("Start application")

	p := NewPrefetchBuffer(ctx, logger)
	p.Start()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()
	sig := <-c
	logger.Info(fmt.Sprintf("Caught signal %v", sig))
	// shutdown other goroutines gracefully
	// close other resources
}
