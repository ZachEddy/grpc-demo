package cmd

import (
	"context"
	"fmt"

	"grpc-demo/grpc-demo-server/pkg/pb"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a new clock on the server",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := NewClient()
		if err != nil {
			PrintErrAndQuit(err)
		}
		res, err := client.CreateClock(context.Background(), &pb.CreateTickerRequest{})
		if err != nil {
			PrintErrAndQuit(err)
		}
		fmt.Printf("Clock '%s' started\n", res.GetName())
	},
}
