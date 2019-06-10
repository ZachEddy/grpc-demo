package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"grpc-demo/grpc-demo-server/pkg/pb"

	"github.com/spf13/cobra"
)

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stops a clock on the server",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := NewClient()
		stream, err := client.GetClockEvents(context.Background(), &pb.GetClockEventsRequest{})
		if err != nil {
			PrintErrAndQuit(err)
		}
		doneChan := make(chan (interface{}))
		go func() {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Press enter to stop streaming")
			if _, err := reader.ReadString('\n'); err != nil {
				PrintErrAndQuit(err)
			}
			if err := stream.CloseSend(); err != nil {
				PrintErrAndQuit(err)
			}
			close(doneChan)
			return
		}()
		go func() {
			for {
				event, err := stream.Recv()
				if err == io.EOF {
					close(doneChan)
					return
				}
				log.Printf("%s: %s", event.Name, event.Event)
			}
		}()
		<-doneChan
	},
}
