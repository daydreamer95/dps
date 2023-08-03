package main

import (
	"context"
	"dps/internal/pkg"
	"dps/internal/pkg/config"
	"dps/logger"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var confFile = flag.String("c", "conf.yaml", "Path to the server configuration file.")

func main() {
	flag.Parse()
	ctx := context.Background()
	config.MustLoadConfig(*confFile)

	ctc := make(chan pkg.Topic)
	d := make(chan pkg.Item)
	var r pkg.IReplenishsesWorker
	r = pkg.NewReplenishesWorker(ctx, ctc, d)
	go r.Start()

	dequeWorker := pkg.NewDequeueWorker(ctx, d)
	go dequeWorker.Start()

	srv := pkg.NewGrpcServer(r)
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
