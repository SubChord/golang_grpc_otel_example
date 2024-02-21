package interceptor

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

func LogrusServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	t0 := time.Now()
	logrus.Infof("Request: %v", req)
	resp, err := handler(ctx, req)
	logrus.Infof("Response: %v (%v)", resp, time.Since(t0))
	return resp, err
}

func LogrusClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	logrus.Infof("Request: %v", req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	logrus.Infof("Response: %v", reply)
	return err
}
