package main

import (
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc"
	"log"
	"pmg-sample/internal/client"
	"pmg-sample/pkg/pb"
	"time"
)

const address = "server:5333"

func main() {
	creds, err := client.LoadTLSCredentials()
	if err != nil {
		log.Fatalf("error with loadTls: %v", err)
	}
	auth := client.AuthCreds{
		Login:    "root",
		Password: "root",
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

	catClient := pb.NewCatalogueClient(conn)
	// bidirectional block
	list := []*pb.Item{
		{Title: "Inteo Four f3-6766x", Description: "4-core, 8-thread mobile processor", Price: 563, Stock: false},
		{Title: "Baikel Wolf", Description: "32-core, 64-thread desktop processor", Price: 872, Stock: true},
		{Title: "Gifox Genvideo", Description: "8gb VRAM mobile GPU", Price: 767, Stock: false},
		{Title: "ATD Vavedon 11800", Description: "6gb VRAM desktop GPU", Price: 344, Stock: true},
	}
	if err = client.AddListItem(catClient, list); err != nil {
		log.Fatalf("failed to add list of items: %v", err)
	}
	// unary block
	for {
		item := &pb.Item{
			Title:       gofakeit.CarModel(),
			Description: gofakeit.CarType(),
			Price:       gofakeit.Price(100, 15000),
			Stock:       gofakeit.Bool(),
		}
		itemId, err := client.AddItem(catClient, item) // add item
		if err != nil {
			log.Fatalf("add item failed: %v", err)
		}
		tableItem, err := client.GetItem(catClient, itemId) // get item by id
		if err != nil {
			log.Fatalf("get item failed: %v", err)
		}
		log.Printf("[unary] item in table: %v", tableItem)
		time.Sleep(5 * time.Second)
	}
}
