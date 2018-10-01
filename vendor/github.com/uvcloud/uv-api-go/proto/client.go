package uv

import (
	"crypto/tls"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//Client wrapper for UV
type Client interface {
	Close()
	V2() UVClient
	Context() context.Context
}
type grpcClient struct {
	Controller string
	Conn       *grpc.ClientConn
	Cancel     context.CancelFunc
	UV         UVClient
	Ctx        context.Context
}

func (client *grpcClient) Close() {
	client.Conn.Close()
	client.Cancel()
}

func (client *grpcClient) V2() UVClient {
	return client.UV
}

func (client *grpcClient) Context() context.Context {
	return client.Ctx
}

//Connect init a connection to grpc server
func Connect(host string, perRPC credentials.PerRPCCredentials) Client {
	client := new(grpcClient)

	client.Controller = host
	opts := []grpc.DialOption{
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		grpc.WithPerRPCCredentials(perRPC),

		// oauth.NewOauthAccess requires the configuration of transport credentials.
		grpc.WithTransportCredentials(
			//TODO skip this for now
			credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
		),
	}
	conn, err := grpc.Dial(client.Controller, opts...)
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	client.Conn = conn
	client.UV = NewUVClient(client.Conn)
	client.Ctx, client.Cancel = context.WithTimeout(context.Background(), 10*time.Second)
	return client
}
