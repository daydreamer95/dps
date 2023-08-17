package pkg

import (
	"dps/internal/pkg/config"
	"dps/internal/pkg/dps_pb"
	"dps/internal/pkg/entity"
	"dps/logger"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type DpsServer struct {
	grpcSrv *grpc.Server
}

func NewGrpcServer(rpw IReplenishsesWorker,
	topicProcessor entity.ITopicProcessor,
	itemProcessor entity.IItemProcessor) *DpsServer {
	out := &DpsServer{}
	var opts []grpc.ServerOption
	out.grpcSrv = grpc.NewServer(opts...)
	dps_pb.RegisterDpsServiceServer(
		out.grpcSrv,
		NewRouterGrpc(rpw, topicProcessor, itemProcessor))
	return out
}

func (d *DpsServer) StartListenAndServer() {
	port := config.Config.App.Port
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Config.App.Host, port))
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
