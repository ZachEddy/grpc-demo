package cmd

import (
	"context"
	"fmt"

	"grpc-demo/grpc-demo-server/pkg/pb"

	"github.com/spf13/cobra"
)

var stopAllCmd = &cobra.Command{
	Use:   "stop-all",
	Short: "Stops a clock on the server",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := NewClient()
		res, err := client.StopAllClocks(context.Background(), &pb.StopAllClocksRequest{})
		PrintErrAndQuit(err)
		for _, name := range res.GetNames() {
			fmt.Println(name)
		}
	},
}
