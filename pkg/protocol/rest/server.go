package rest

import (
	v1 "axxon/pkg/api/v1"
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func RunServer(ctx context.Context, grpcPort string, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1.RegisterFetchServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		log.Fatalf("Failed to start grpc gateway: %v", err)
	}

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}
	return srv.ListenAndServe()
}
