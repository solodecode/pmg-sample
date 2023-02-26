package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"net"
	"os"
	"pmg-sample/pb"
)

const (
	dbUrl = "postgres://yourlogin:yourpass@localhost:5432/yourdb"
)

type server struct {
	pb.UnimplementedCatalogueServer           // must be implemented
	conn                            *pgx.Conn // db connection
}

func (s *server) GetItem(ctx context.Context, id *wrapperspb.UInt64Value) (*pb.Item, error) {
	resp := &pb.Item{}
	if err := s.conn.QueryRow(ctx, fmt.Sprintf("select * from test WHERE id=%d", id.Value)).Scan(&resp.Id, &resp.Title, &resp.Description, &resp.Price, &resp.Stock); err != nil {
		log.Printf("QueryRow failed: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *server) AddItem(ctx context.Context, item *pb.Item) (*wrapperspb.UInt64Value, error) {
	var id wrapperspb.UInt64Value
	if err := s.conn.QueryRow(ctx, fmt.Sprintf("insert into test (title, description, price, stock) values ('%s', '%s', %f, %t) RETURNING id", item.Title, item.Description, item.Price, item.Stock)).Scan(&id.Value); err != nil {
		log.Printf("QueryRow failed: %v", err)
		return nil, err
	}
	return &id, nil
}

func LoadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("crt/server/server.crt", "crt/server/server.key")
	if err != nil {
		return nil, err
	}

	ca, err := os.ReadFile("crt/ca/ca.crt")
	if err != nil {
		return nil, err
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(ca) {
		return nil, errors.New("cannot load ca-crt to pool")
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    capool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer conn.Close(context.Background())

	lis, err := net.Listen("tcp", ":53452")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := LoadTLSCredentials()
	if err != nil {
		log.Fatalf("cannot load tls certs: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterCatalogueServer(s, &server{conn: conn})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
