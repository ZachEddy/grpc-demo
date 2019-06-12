package cmd

import (
	"fmt"
	"os"

	"grpc-demo/demo-server/pkg/pb"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(stopAllCmd)
	rootCmd.AddCommand(streamCmd)
}

var rootCmd = &cobra.Command{
	Short: "This is an awesome demo CLI client for the gRPC server",
}

// Execute launches the root command
func Execute() error {
	return rootCmd.Execute()
}

func printErrAndQuit(err error) {
	fmt.Printf("unable to complete request: %v\n", err)
	os.Exit(0)
}

func newClient() (pb.DemoClient, error) {
	// TODO: hardcoded, but whatever
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewDemoClient(conn), nil
}
