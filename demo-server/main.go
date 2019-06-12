package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"grpc-demo/demo-server/pkg/demo"
	"grpc-demo/demo-server/pkg/pb"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	flagServerPort  = flag.Int("server-address", 8080, "grpc server port")
	flagMetricsPort = flag.Int("metrics-address", 7070, "metrics port")
)

func init() {
	flag.Parse()
}

func main() {
	log.Printf("serving grpc on port %d\n", *flagServerPort)
	log.Printf("serving prometheus metrics on port %d\n", *flagMetricsPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to create tcp listener: %v", err)
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	)
	pb.RegisterDemoServer(server, demo.NewDemoServer())
	grpc_prometheus.Register(server)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":7070", nil)
	}()
	server.Serve(lis)
}
