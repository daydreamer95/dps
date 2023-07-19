package dps_srv

import (
	"dps/logger"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type DpsServer struct {
	grpcSrv *grpc.Server
}

func NewGrpcServer() *DpsServer {
	out := &DpsServer{}
	var opts []grpc.ServerOption
	out.grpcSrv = grpc.NewServer(opts...)
	RegisterDpsServiceServer(out.grpcSrv, NewRouterGrpc())
	return out
}

func (d *DpsServer) StartListenAndServer() {
	port := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to listen: %v", err))
	}
	logger.Info(fmt.Sprintf("failed to listen: %v", port))
	d.grpcSrv.Serve(lis)
}
