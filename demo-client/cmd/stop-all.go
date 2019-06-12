package cmd

import (
	"context"
	"fmt"

	"grpc-demo/demo-server/pkg/pb"

	"github.com/spf13/cobra"
)

var stopAllCmd = &cobra.Command{
	Use:   "stop-all",
	Short: "Stops all the the server clocks",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := newClient()
		res, err := client.StopAllClocks(context.Background(), &pb.StopAllClocksRequest{})
		if err != nil {
			printErrAndQuit(err)
		}
		for _, name := range res.GetNames() {
			fmt.Println(name)
		}
	},
}
