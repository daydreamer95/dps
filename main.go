package main

import (
	"context"
	"dps/internal/pkg"
	"dps/internal/pkg/config"
	"dps/internal/pkg/entity"
	"dps/internal/pkg/storage/registry"
	"dps/logger"
	"flag"
	"fmt"
	"go.uber.org/zap"
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
	registry.MustWaitStoreUp()

	ctc := make(chan entity.Topic)
	d := make(chan entity.Item)
	expChan := make(chan entity.Item)
	var r pkg.IReplenishsesWorker
	r = pkg.NewReplenishesWorker(ctx, ctc, expChan, d)
	go r.Start()

	dequeWorker := pkg.NewDequeueWorker(ctx, d)
	go dequeWorker.Start()

	topicProcessor := entity.NewTopicProcessor()
	itemProcessor := entity.NewItemProcessor()
	srv := pkg.NewGrpcServer(r, topicProcessor, itemProcessor)
	go srv.StartListenAndServer()

	// process lease message
	go func(processor entity.IItemProcessor) {
		for {
			select {
			case expItem := <-expChan:
				err := itemProcessor.Delete(context.TODO(), expItem.TopicId, expItem.Id)
				if err != nil {
					logger.Info("Delete fail expire item id [%v]", zap.String("item_id", expItem.Id))
					return
				}
				logger.Info("Delete expire message", zap.String("item_id", expItem.Id), zap.Int64("lease_after", expItem.LeaseAfter.Unix()))
			}
		}
	}(itemProcessor)

	//Pprof http
	go func() {
		log.Print(http.ListenAndServe("localhost:6060", nil))
	}()

	// application handle panic
	go handlePanic()

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

func handlePanic() {
	if r := recover(); r != nil {
		logger.Error(fmt.Sprintf("Panic occur, detail: %v", r))
	}
}
