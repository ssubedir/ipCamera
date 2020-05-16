package service

import (
	"context"
	"fmt"

	proto "ipCamera/proto/ipcamera/proto"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Stream(ctx context.Context, req *proto.StreamRequest) (*proto.StreamResponse, error) {

	fmt.Println("Client-", req.Data)
	return &proto.StreamResponse{Response: "We got your message"}, nil
}
