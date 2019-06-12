package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"grpc-demo/demo-server/pkg/pb"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops a server clock",
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
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Choose clock to stop: ")
		responseRaw, err := reader.ReadString('\n')
		if err != nil {
			printErrAndQuit(err)
		}
		response, err := strconv.Atoi(strings.TrimSpace(responseRaw))
		if err != nil {
			printErrAndQuit(err)
		}
		if response < 0 || response > len(clocks.GetNames())-1 {
			printErrAndQuit(fmt.Errorf("Invalid clock choice"))
		}
		if _, err := client.StopClock(context.Background(), &pb.StopClockRequest{
			Name: clocks.GetNames()[response],
		}); err != nil {
			printErrAndQuit(err)
		}
	},
}
