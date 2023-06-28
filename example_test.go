package grpcmtls_test

import (
	context "context"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/covrom/grpcmtls"
	"github.com/covrom/grpcmtls/api"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func ExampleClientServerGrpc() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("starting...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	// certificates
	certs := grpcmtls.CertFiles{
		ServerCertPemFilePath: "./certs/server.pem",
		ServerKeyPemFilePath:  "./certs/server.key",
		CAFilePath:            "./certs/cacert.pem",
		ClientCertPemFilePath: "./certs/client.pem",
		ClientKeyPemFilePath:  "./certs/client.key",
	}

	wg := &sync.WaitGroup{}

	// serve grpc
	wg.Add(1)
	go ServeEcho(ctx, "localhost:2356", certs.ServerCertFiles(), wg)

	// pause
	time.Sleep(100 * time.Millisecond)

	// connect client
	cli, err := grpcmtls.NewGrpcClient(ctx, "localhost:2356", certs.ClientCertFiles(), api.NewEchoServiceClient)
	if err != nil {
		slog.Error("NewGrpcClient", "err", err)
	} else {
		res, err := cli.Echo(ctx, &api.Source{Msg: "hello"})
		if err != nil {
			slog.Error("cli.Echo", "err", err)
		} else {
			slog.Info("echo result", "msg", res.Msg)
		}
	}

	// finalize
	cancel()
	wg.Wait()

	slog.Info("exit.")
}

func ServeEcho(ctx context.Context, addr string,
	certs grpcmtls.ServerCertFiles, wg *sync.WaitGroup) {
	defer wg.Done()

	srvc := api.NewEchoGrpcService()

	l, e := net.Listen("tcp", addr)
	if e != nil {
		slog.Error("GRPC listen error: ", "err", e)
		return
	}

	grpcServer, err := grpcmtls.NewServerGRPC(certs,
		func(srv *grpc.Server) {
			api.RegisterEchoServiceServer(srv, srvc)
		})

	if err != nil {
		slog.Error("NewServerGRPC", "err", err)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("GRPC server listen", "addr", addr)
		if err := grpcServer.Serve(l); err != nil {
			slog.Error("failed to serve grpc", "err", err)
		}
		slog.Info("GRPC server stopped")
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
}
