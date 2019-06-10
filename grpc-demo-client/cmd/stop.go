package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"grpc-demo/grpc-demo-server/pkg/pb"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops a clock on the server",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := NewClient()
		if err != nil {
			PrintErrAndQuit(err)
		}
		clocks, err := client.ListClocks(context.Background(), &pb.ListClocksRequest{})
		if err != nil {
			PrintErrAndQuit(err)
		}
		for i, clock := range clocks.GetNames() {
			fmt.Printf("%d: %s\n", i, clock)
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Choose clock to stop: ")
		responseRaw, err := reader.ReadString('\n')
		if err != nil {
			PrintErrAndQuit(err)
		}
		response, err := strconv.Atoi(strings.TrimSpace(responseRaw))
		if err != nil {
			PrintErrAndQuit(err)
		}
		if response < 0 || response > len(clocks.GetNames())-1 {
			PrintErrAndQuit(fmt.Errorf("fuck you"))
		}
		if _, err := client.StopClock(context.Background(), &pb.StopClockRequest{
			Name: clocks.GetNames()[response],
		}); err != nil {
			PrintErrAndQuit(err)
		}
	},
}
