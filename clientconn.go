package grpcmtls

import (
	context "context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math"
	"os"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewClientConnGrpc(ctx context.Context, addr string, certs ClientCertFiles) (*grpc.ClientConn, error) {
	cert, err := tls.LoadX509KeyPair(certs.ClientCertPemFilePath, certs.ClientKeyPemFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert: %w", err)
	}

	ca := x509.NewCertPool()
	caBytes, err := os.ReadFile(certs.CAFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ca cert %q: %w", certs.CAFilePath, err)
	}

	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		return nil, fmt.Errorf("failed to parse %q", certs.CAFilePath)
	}

	tlsConfig := &tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	opts := []grpcretry.CallOption{
		grpcretry.WithBackoff(grpcretry.BackoffExponential(100 * time.Millisecond)),
		grpcretry.WithMax(5),
	}

	maxMessageSize := math.MaxInt32 // 2Gb

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMessageSize),
			grpc.MaxCallSendMsgSize(maxMessageSize),
		),
		grpc.WithStreamInterceptor(
			grpcmiddleware.ChainStreamClient(
				grpcretry.StreamClientInterceptor(opts...),
			),
		),
		grpc.WithUnaryInterceptor(
			grpcmiddleware.ChainUnaryClient(
				grpcretry.UnaryClientInterceptor(opts...),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
