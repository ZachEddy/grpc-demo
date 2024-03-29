package cmd

import (
	"context"
	"fmt"

	"grpc-demo/demo-server/pkg/pb"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a new server clock",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := newClient()
		if err != nil {
			printErrAndQuit(err)
		}
		res, err := client.CreateClock(context.Background(), &pb.CreateTickerRequest{})
		if err != nil {
			printErrAndQuit(err)
		}
		fmt.Printf("Clock '%s' started\n", res.GetName())
	},
}
