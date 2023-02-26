package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
	"pmg-sample/pb"
)

const address = "127.0.0.1:53452"

func main() {
	creds, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf("error with loadTls: %v", err)
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v")
	}
	defer conn.Close()
	client := pb.NewCatalogueClient(conn)

	itemId, err := client.AddItem(context.Background(), &pb.Item{
		Title:       "ATD Ruzen 16",
		Description: "16-core, 32-thread desktop processor",
		Price:       564,
		Stock:       true,
	})
	if err != nil {
		log.Fatalf("add item failed: %v", err)
	}

	item, err := client.GetItem(context.Background(), itemId)
	if err != nil {
		log.Fatalf("get item failed: %v", err)
	}
	fmt.Println(item)
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	certs, err := tls.LoadX509KeyPair("crt/client/client.crt", "crt/client/client.key")
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
		Certificates: []tls.Certificate{certs},
		RootCAs:      capool,
	}
	return credentials.NewTLS(config), nil
}
