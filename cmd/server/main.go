package main

import (
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"pmg-sample/internal/server"
	"pmg-sample/pkg/pb"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer conn.Close(context.Background())

	lis, err := net.Listen("tcp", ":5333")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := server.LoadTLSCredentials()
	if err != nil {
		log.Fatalf("cannot load tls certs: %v", err)
	}

	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(server.UnaryAuthInterceptor),
		grpc.StreamInterceptor(server.ServerStreamInterceptor),
	)
	pb.RegisterCatalogueServer(s, &server.Server{Conn: conn})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
