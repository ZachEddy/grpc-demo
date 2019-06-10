package cmd

import (
	"fmt"
	"os"

	"grpc-demo/demo-server/pkg/pb"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(stopAllCmd)
	rootCmd.AddCommand(streamCmd)
}

var rootCmd = &cobra.Command{
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
}

func Execute() error {
	return rootCmd.Execute()
}

func PrintErrAndQuit(err error) {
	fmt.Printf("unable to complete request: %v\n", err)
	os.Exit(0)
}

func NewClient() (pb.DemoClient, error) {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewDemoClient(conn), nil
}
