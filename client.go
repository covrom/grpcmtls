package grpcmtls

import (
	context "context"

	"google.golang.org/grpc"
)

func NewGrpcClient[T any](ctx context.Context,
	addr string,
	certs ClientCertFiles,
	constructor func(cc grpc.ClientConnInterface) T,
) (T, error) {
	conn, err := NewClientConnGrpc(ctx, addr, certs)
	if err != nil {
		return *(new(T)), err
	}

	return constructor(conn), nil
}
