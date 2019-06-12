package cmd

import (
	"context"
	"fmt"

	"grpc-demo/demo-server/pkg/pb"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all server clocks",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := newClient()
		if err != nil {
			printErrAndQuit(err)
		}
		clocks, err := client.ListClocks(context.Background(), &pb.ListClocksRequest{})
		if err != nil {
			printErrAndQuit(err)
		}
		for i, clock := range clocks.GetNames() {
			fmt.Printf("%d: %s\n", i, clock)
		}
	},
}
