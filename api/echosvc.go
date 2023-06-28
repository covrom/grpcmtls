package api

import (
	context "context"
)

type EchoGrpcService struct {
	UnimplementedEchoServiceServer
}

func NewEchoGrpcService() *EchoGrpcService {
	return &EchoGrpcService{}
}

func (s *EchoGrpcService) Echo(ctx context.Context, src *Source) (*Result, error) {
	return &Result{Msg: src.GetMsg()}, nil
}
