package common

import (
	"github.com/kurojs/consul-grpc/consul"
	"google.golang.org/grpc"
)

// GrpcClient ...
type GrpcClient struct{}

//NewGrpcClient ...
func NewGrpcClient(address string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithInsecure())

	return grpc.Dial(address, opts...)
}

// NewGrpcClientWithLoadBalance ...
func NewGrpcClientWithLoadBalance(serviceName, tag, address string) (*grpc.ClientConn, error) {
	resolver, err := consul.NewResolver(serviceName, tag)
	if err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBalancer(grpc.RoundRobin(resolver)))

	return grpc.Dial(address, opts...)
}
