package main

import (
	"context"
	"dps/internal/pkg"
	"dps/internal/pkg/config"
	"dps/internal/pkg/entity"
	"dps/logger"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

var confFile = flag.String("c", "conf.yaml", "Path to the server configuration file.")

func main() {
	flag.Parse()
	ctx := context.Background()
	config.MustLoadConfig(*confFile)

	ctc := make(chan entity.Topic)
	d := make(chan entity.Item)
	var r pkg.IReplenishsesWorker
	r = pkg.NewReplenishesWorker(ctx, ctc, d)
	go r.Start()

	dequeWorker := pkg.NewDequeueWorker(ctx, d)
	go dequeWorker.Start()

	topicProcessor := entity.NewTopicProcessor()
	itemProcessor := entity.NewItemProcessor()
	srv := pkg.NewGrpcServer(r, topicProcessor, itemProcessor)
	go srv.StartListenAndServer()

	go func() {
		log.Print(http.ListenAndServe("localhost:6060", nil))
	}()
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
