package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.360live.vn/zpi/consul/common"
	greeting "gitlab.360live.vn/zpi/consul/grpc-gen"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func Test_HealthCheck(t *testing.T) {
	clientConn, err := common.NewGrpcClient("localhost:8001")
	require.Nil(t, err)

	res, err := greeting.NewGreetingClient(clientConn).SayHi(context.Background(), &greeting.Request{})
	require.Nil(t, err)
	fmt.Println(res)

	resp, err := grpc_health_v1.NewHealthClient(clientConn).Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	assert.Nil(t, err)
	fmt.Println(resp)
}
