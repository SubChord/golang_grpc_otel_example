package main

import (
	"bitbucket.org/be-mobile/jaeger-tracing-lib/tracer"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc_example/pkg/interceptor"
	"grpc_example/pkg/message/v1"

	"net"
)

func main() {
	appContext := context.Background()

	// init OpenTelemetry exporter
	exporter, err := tracer.InitGlobalDebugTracer("grpc_example_server", "")
	if err := err; err != nil {
		logrus.Fatal(err)
	}
	defer exporter.Shutdown(appContext)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LogrusServerInterceptor),
		grpc.StatsHandler(otelgrpc.NewServerHandler()), // add OpenTelemetry stats handler
	)
	reflection.Register(server)

	message.RegisterMessageServiceServer(server, MessageServiceServer{})

	listener, err := net.Listen("tcp", ":8088")
	if err != nil {
		panic(err)
	}

	go func() {
		<-appContext.Done()
		server.GracefulStop()
	}()

	logrus.Info("Starting gRPC server")
	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}

type MessageServiceServer struct {
	message.UnimplementedMessageServiceServer
}

func (m MessageServiceServer) Send(ctx context.Context, request *message.MessageRequest) (*message.MessageResponse, error) {
	ctx, span := otel.Tracer("MessageServiceServer").Start(ctx, "Send")
	defer span.End()

	span.AddEvent("Received message")

	return &message.MessageResponse{Response: fmt.Sprintf("echo: %v", request.Message)}, nil
}
