package dps_srv

import (
	"dps/internal/pkg"
	"dps/internal/pkg/logger"
	"dps/internal/pkg/repository"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type DpsServer struct {
	grpcSrv *grpc.Server
}

func NewGrpcServer(rpw pkg.IReplenishsesWorker,
	itemRepository *repository.ItemRepository,
	topicRepository *repository.TopicRepository) *DpsServer {
	out := &DpsServer{}
	var opts []grpc.ServerOption
	out.grpcSrv = grpc.NewServer(opts...)
	RegisterDpsServiceServer(out.grpcSrv, NewRouterGrpc(rpw, itemRepository, topicRepository))
	return out
}

func (d *DpsServer) StartListenAndServer() {
	port := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to listen: %v casue err [%v]", port, err))
	}
	go func() {
		err := d.grpcSrv.Serve(lis)
		if err != nil {
			logger.Fatal(fmt.Sprintf("Error start grpc detail: [%v]", err))
		}
	}()
	logger.Info(fmt.Sprintf("Success start grpc on port [%v]", port))
}
