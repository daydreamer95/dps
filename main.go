package main

import (
	"context"
	"dps/dps_srv"
	"dps/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	logger.Info("Start application")

	p := NewPrefetchBuffer(ctx)
	p.Start()

	srv := dps_srv.NewGrpcServer()
	go srv.StartListenAndServer()

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
