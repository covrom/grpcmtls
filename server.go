package grpcmtls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math"
	"os"
	"runtime"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func NewServerGRPC(certs ServerCertFiles, reg func(srv *grpc.Server), opt ...grpc.ServerOption) (*grpc.Server, error) {
	maxMessageSize := math.MaxInt32 // 2Gb
	recoveryOption := grpcrecovery.WithRecoveryHandlerContext(panicHandler())

	cert, err := tls.LoadX509KeyPair(certs.ServerCertPemFilePath, certs.ServerKeyPemFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair: %w", err)
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
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
		RootCAs:      ca,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	s := grpc.NewServer(append(opt,
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.MaxRecvMsgSize(maxMessageSize),
		grpc.MaxSendMsgSize(maxMessageSize),
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcrecovery.StreamServerInterceptor(recoveryOption)),
		),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcrecovery.UnaryServerInterceptor(recoveryOption),
		)),
	)...)

	reg(s)

	// reflection.Register(s)
	return s, nil
}

func panicHandler() func(ctx context.Context, p interface{}) error {
	return func(ctx context.Context, p interface{}) error {
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		slog.Error("panic", "message", p, "stack", string(buf[0:stackSize]))
		return status.Error(codes.Internal, "internal error")
	}
}
