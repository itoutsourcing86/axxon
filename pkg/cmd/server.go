package cmd

import (
	"axxon/pkg/protocol/grpc"
	"axxon/pkg/protocol/rest"
	v1 "axxon/pkg/service/v1"
	"context"
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	GRPCPort   string
	HTTPPort   string
	DBHost     string
	DBUser     string
	DBPassword string
	// Database name
	DBSchema string
}

func RunServer() error {
	var cfg Config

	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "GRPC port")
	flag.StringVar(&cfg.HTTPPort, "http-port", "", "HTTP port")
	flag.StringVar(&cfg.DBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DBUser, "db-user", "", "Database username")
	flag.StringVar(&cfg.DBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("Invalid GRPC port: %s", cfg.GRPCPort)
	}

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid HTTP port: '%s'", cfg.HTTPPort)
	}

	ctx := context.Background()
	dataStoreName := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBSchema,
	)
	// fmt.Println(dataStoreName)

	db, err := sql.Open("mysql", dataStoreName)
	if err != nil {
		return err
	}
	defer db.Close()

	v1API := v1.NewFetchServiceServer(db)

	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
