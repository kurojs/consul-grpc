package server

import (
	"context"
	"log"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"gitlab.360live.vn/zpi/consul/common"
	"gitlab.360live.vn/zpi/consul/consul"
	greeting "gitlab.360live.vn/zpi/consul/grpc-gen"
	"google.golang.org/grpc"
)

//Server ...
type Server struct {
	name     string
	port     int
	hostname string
	id       string
	tags     []string
}

//NewServer ...
func NewServer(name, hostname string, tags []string, port int) *Server {
	return &Server{
		name:     name,
		port:     port,
		hostname: hostname,
		tags:     tags,
	}
}

//SayHi ...
func (srv *Server) SayHi(ctx context.Context, in *greeting.Request) (*greeting.Response, error) {
	resp := &greeting.Response{
		Data: "Received " + in.GetData() + "port: " + strconv.Itoa(srv.port) + ", id: " + srv.id,
		Time: time.Now().Unix(),
		Id:   strconv.Itoa(srv.port),
	}

	return resp, nil
}

//RegisterServer ...
func (srv *Server) registerServer(server *grpc.Server) {
	greeting.RegisterGreetingServer(server, srv)
}

//Run ...
func (srv *Server) Run() error {
	srv.id = srv.name + "-" + uuid.Must(uuid.NewV4()).String()
	grpcSrv := common.NewGRPCServer(srv.registerServer)
	consultService, err := consul.NewService(srv.id, srv.name, srv.hostname, srv.port, srv.tags)
	if err != nil {
		log.Fatalln("Create consul server failed", err)
		return err
	}

	grpcSrv.WithConsul(consultService)
	err = grpcSrv.RegisterConsul()
	if err != nil {
		log.Fatalln("Register consul failed", err)
		return err
	}

	grpcSrv.RegisterHealthCheck(&HealthServer{})

	err = grpcSrv.Run(srv.port)
	if err != nil {
		log.Fatalln("Create gRPC server failed", err)
		return err
	}

	return nil
}
