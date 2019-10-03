package common

import (
	"gitlab.360live.vn/zpi/consul/consul"
	"google.golang.org/grpc"
)

type GrpcClient struct{}

//NewGrpcClient ...
func NewGrpcClient(address string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithInsecure())

	return grpc.Dial(address, opts...)
}

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
