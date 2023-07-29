package main

import (
	"context"
	"dps/internal/pkg"
	"dps/internal/pkg/dps_srv"
	"dps/internal/pkg/logger"
	"dps/internal/pkg/repository"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal(fmt.Sprint("Error load env file: ", err))
	}

	db, err := pkg.NewMySqlDb()
	if err != nil {
		logger.Fatal(fmt.Sprint("Error connect to db: ", err))
	}
	itemRepo := repository.NewItemRepository(db)
	topicRepo := repository.NewTopicRepository(db)

	ctc := make(chan string)
	d := make(chan pkg.Item)
	var r pkg.IReplenishsesWorker
	r = pkg.NewReplenishesWorker(ctx, ctc, d)
	go r.Start()

	srv := dps_srv.NewGrpcServer(r, itemRepo, topicRepo)
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
