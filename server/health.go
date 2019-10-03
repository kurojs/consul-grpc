package server

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

//HealthServer ...
type HealthServer struct {
	port     int
	hostname string
	id       string
}

//Check ...
func (srv *HealthServer) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

//Watch ...
func (srv *HealthServer) Watch(in *grpc_health_v1.HealthCheckRequest, ws grpc_health_v1.Health_WatchServer) error {
	ws.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})

	return nil
}
