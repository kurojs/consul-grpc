package common

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"gitlab.360live.vn/zpi/consul/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//GrpcService re-implement common grpc-service
type GrpcService struct {
	server       *grpc.Server
	port         int
	listener     net.Listener
	grpcRegister GrpcRegister
	hooks        []HookFunc
	consul       *consul.Service
	done         chan error
}

//GrpcRegister ...
type GrpcRegister func(server *grpc.Server)

//HookFunc ...
type HookFunc func()

//NewGRPCServer ...
func NewGRPCServer(register GrpcRegister) *GrpcService {
	return &GrpcService{
		grpcRegister: register,
		server:       grpc.NewServer(),
	}
}

//Run ...
func (s *GrpcService) Run(port int) error {
	var err error
	s.port = port
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	s.grpcRegister(s.server)

	go func() {
		fmt.Println("Server is running on port:", s.port)
		err = s.server.Serve(s.listener)
	}()

	sigs := make(chan os.Signal, 1)
	s.done = make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		s.runHook()
		err := s.DeregisterConsul()
		err = s.Shutdown()
		s.server.GracefulStop()
		s.done <- err
	}()
	err = <-s.done
	return err
}

func (s *GrpcService) runHook() {
	for _, hook := range s.hooks {
		defer hook()
	}
}

// Shutdown -
func (s *GrpcService) Shutdown() error {
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return err
		}
		s.listener = nil
	}
	return nil
}

// AddShutdownHook ...
func (s *GrpcService) AddShutdownHook(hookFunc HookFunc) {
	s.hooks = append(s.hooks, hookFunc)
}

//WithConsul ...
func (s *GrpcService) WithConsul(consulService *consul.Service) {
	s.consul = consulService
}

//RegisterConsul ...
func (s *GrpcService) RegisterConsul() error {
	if s.consul != nil {
		return s.consul.Register()
	}

	return nil
}

//DeregisterConsul ...
func (s *GrpcService) DeregisterConsul() error {
	if s.consul != nil {
		return s.consul.Deregister()
	}

	return nil
}

//RegisterHealthCheck ...
func (s *GrpcService) RegisterHealthCheck(server HealthCheck) {
	grpc_health_v1.RegisterHealthServer(s.server, server)
}

//HealthCheck ...
type HealthCheck interface {
	Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error)
	Watch(in *grpc_health_v1.HealthCheckRequest, ws grpc_health_v1.Health_WatchServer) error
}
