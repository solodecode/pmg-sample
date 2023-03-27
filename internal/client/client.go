package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"log"
	"os"
	"pmg-sample/pkg/pb"
)

type AuthCreds struct {
	Login    string
	Password string
}

func (auth AuthCreds) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	creds := auth.Login + ":" + auth.Password
	enc := base64.StdEncoding.EncodeToString([]byte(creds))
	return map[string]string{
		"authorization": enc,
	}, nil
}

func (auth AuthCreds) RequireTransportSecurity() bool {
	return true
}

func AddListItem(client pb.CatalogueClient, items []*pb.Item) error {
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

func GetItem(client pb.CatalogueClient, itemId *wrapperspb.UInt64Value) (*pb.Item, error) {
	item, err := client.GetItem(context.Background(), itemId)
	if err != nil {
		return item, err
	}
	return item, nil
}

func AddItem(client pb.CatalogueClient, item *pb.Item) (*wrapperspb.UInt64Value, error) {
	itemId, err := client.AddItem(context.Background(), item)
	if err != nil {
		return itemId, err
	}
	return itemId, nil
}

func LoadTLSCredentials() (credentials.TransportCredentials, error) {
	certs, err := tls.LoadX509KeyPair("config/crt/client/client.crt", "config/crt/client/client.key")
	if err != nil {
		return nil, err
	}

	ca, err := os.ReadFile("config/crt/ca/ca.crt")
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
