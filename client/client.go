package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"log"
	"os"
	"pmg-sample/pb"
)

const address = "127.0.0.1:53452"

type authCreds struct {
	login    string
	password string
}

func (auth authCreds) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	creds := auth.login + ":" + auth.password
	enc := base64.StdEncoding.EncodeToString([]byte(creds))
	return map[string]string{
		"authorization": enc,
	}, nil
}

func (auth authCreds) RequireTransportSecurity() bool {
	return true
}

func addListItem(client pb.CatalogueClient, items []*pb.Item) error {
	stream, err := client.AddListItem(context.Background())
	if err != nil {
		return err
	}
	// send items to server
	for _, v := range items {
		if err = stream.Send(v); err != nil {
			return err
		}
	}
	// close send
	if err = stream.CloseSend(); err != nil {
		return err
	}
	// receiving ID's from server
	for {
		id, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		log.Printf("[bidirectional] product id: %d", id.Value)
	}
	return nil
}

func getItem(client pb.CatalogueClient, itemId *wrapperspb.UInt64Value) (*pb.Item, error) {
	item, err := client.GetItem(context.Background(), itemId)
	if err != nil {
		return item, err
	}
	return item, nil
}

func addItem(client pb.CatalogueClient, item *pb.Item) (*wrapperspb.UInt64Value, error) {
	itemId, err := client.AddItem(context.Background(), item)
	if err != nil {
		return itemId, err
	}
	return itemId, nil
}

func main() {
	creds, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf("error with loadTls: %v", err)
	}
	auth := authCreds{
		login:    "root",
		password: "root",
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds), // mTLS security
		grpc.WithPerRPCCredentials(auth),     // Simple authentication
	}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v")
	}
	defer conn.Close()

	client := pb.NewCatalogueClient(conn)

	item := &pb.Item{
		Title:       "ATD Ruzen 16",
		Description: "16-core, 32-thread desktop processor",
		Price:       564,
		Stock:       true,
	}
	// unary block
	itemId, err := addItem(client, item) // add item
	if err != nil {
		log.Fatalf("add item failed: %v", err)
	}
	tableItem, err := getItem(client, itemId) // get item by id
	if err != nil {
		log.Fatalf("get item failed: %v", err)
	}
	log.Printf("[unary] item in table: %v", tableItem)
	// bidirectional block
	list := []*pb.Item{
		{Title: "Inteo Four f3-6766x", Description: "4-core, 8-thread mobile processor", Price: 563, Stock: false},
		{Title: "Baikel Wolf", Description: "32-core, 64-thread desktop processor", Price: 872, Stock: true},
		{Title: "Gifox Genvideo", Description: "8gb VRAM mobile GPU", Price: 767, Stock: false},
		{Title: "ATD Vavedon 11800", Description: "6gb VRAM desktop GPU", Price: 344, Stock: true},
	}
	if err = addListItem(client, list); err != nil {
		log.Fatalf("failed to add list of items: %v", err)
	}
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
