package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"gitlab.360live.vn/zpi/consul/server"
)

var (
	serviceName = "Greeting"
	hostname    = "localhost"
	tags        = []string{"grpc", "consul"}
)

var rootCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		for _, portStr := range args {
			port, err := strconv.Atoi(portStr)
			if err == nil {
				srv := server.NewServer(serviceName, hostname, tags, port)
				go srv.Run()
			}
		}
		fmt.Scanln()
	},
}

//Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
