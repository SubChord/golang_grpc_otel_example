package main

import (
	"bitbucket.org/be-mobile/jaeger-tracing-lib/tracer"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc_example/pkg/interceptor"
	"grpc_example/pkg/message/v1"
	"os"
	"os/signal"
	"time"
)

func main() {
	appContext := context.Background()

	// init OpenTelemetry exporter
	exporter, err := tracer.InitGlobalDebugTracer("grpc_example_client", "")
	if err := err; err != nil {
		logrus.Fatal(err)
	}
	defer exporter.Shutdown(appContext)

	conn, err := grpc.DialContext(
		appContext,
		"localhost:8088",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.LogrusClientInterceptor),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // add OpenTelemetry stats handler
	)
	if err := err; err != nil {
		logrus.Fatal(err)
	}
	defer conn.Close()

	// Start new open telemetry span
	ctx, span := otel.Tracer("main").Start(appContext, "send-message")
	span.AddEvent("Sending message")                                                                                                         // add event to span
	span.SetAttributes(attribute.KeyValue{Key: "timeout", Value: attribute.StringValue(fmt.Sprintf("%v", time.Now().Format(time.RFC3339)))}) // add attribute to span

	_, err = message.NewMessageServiceClient(conn).Send(ctx, &message.MessageRequest{
		Message: "Hello, World!",
	})

	span.End() // end span
	logrus.Info("Message sent")
	if err := err; err != nil {
		logrus.Fatal(err)
	}

	// graceful shutdown
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt)

	<-sigInt
}
