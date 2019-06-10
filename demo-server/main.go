package main

import (
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

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
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
