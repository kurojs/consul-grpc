package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/kurojs/consul-grpc/common"
	greeting "github.com/kurojs/consul-grpc/grpc-gen"
)

func Test_Loadbalance(t *testing.T) {
	conn, err := common.NewGrpcClientWithLoadBalance("Greeting", "grpc", "consul://127.0.0.1:8500")

	require.Nil(t, err)

	client := greeting.NewGreetingClient(conn)

	ctx := context.Background()
	req := greeting.Request{Data: time.Now().String()}

	for i := 0; i < 100; i++ {
		resp, err := client.SayHi(ctx, &req)
		require.Nil(t, err)
		fmt.Println(resp.GetId())
	}
}
