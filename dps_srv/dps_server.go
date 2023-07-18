package dps_srv

import "context"

type DpsServer struct {
	UnimplementedDpsServiceServer
}

func (d *DpsServer) Publish(context.Context, *PublishReq) (*PublishRes, error) {
	return nil, nil
}

func (d *DpsServer) CreateTopic(context.Context, *CreateTopicReq) (*CommonRes, error) {
	return nil, nil
}
func (d *DpsServer) Dequeue(context.Context, *DequeueReq) (*DequeueRes, error) {
	return nil, nil
}
