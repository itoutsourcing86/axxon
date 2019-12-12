package grpc

import (
	v1 "axxon/pkg/api/v1"
	"context"
	"net"

	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, v1API v1.FetchServiceServer, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	v1.RegisterFetchServiceServer(server, v1API)

	return server.Serve(lis)
}
