package main

import (
	"context"
	"dps/internal/pkg"
	"dps/internal/pkg/dps_pb"
	"dps/internal/pkg/entity"
	"dps/logger"
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

	_, err = pkg.NewMySqlDb()
	if err != nil {
		logger.Fatal(fmt.Sprint("Error connect to db: ", err))
	}

	ctc := make(chan string)
	d := make(chan entity.Item)
	var r pkg.IReplenishsesWorker
	r = pkg.NewReplenishesWorker(ctx, ctc, d)
	go r.Start()

	srv := dps_pb.NewGrpcServer(r)
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
