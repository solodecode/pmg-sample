package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"net"
	"os"
	"pmg-sample/pb"
)

const (
	dbUrl = "postgres://yourlogin:yourpass@localhost:5432/yourdb"
)

var (
	errMissMetadata = status.Errorf(codes.InvalidArgument, "missing auth token")
	errLoadCert     = status.Errorf(codes.FailedPrecondition, "cannot load ca-crt to pool")
	errInvalidCreds = status.Errorf(codes.InvalidArgument, "invalid auth credentials")
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

func authorization(header []string) bool {
	if header[0] == base64.StdEncoding.EncodeToString([]byte("root:root")) {
		return true
	}
	return false
}

func unaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissMetadata
	}
	if valid := authorization(md["authorization"]); valid {
		return handler(ctx, req)
	}
	return nil, errInvalidCreds
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
		return nil, errLoadCert
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

	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(unaryAuthInterceptor))
	pb.RegisterCatalogueServer(s, &server{conn: conn})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
