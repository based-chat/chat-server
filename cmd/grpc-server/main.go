// Package main provides the entry point for the chat server.
package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/based-chat/chat-server/utilites"

	srv "github.com/based-chat/chat-server/pkg/chat/v1"
)

const (
	grpcPort = ":50051"
	host     = "0.0.0.0"
)

type server struct {
	srv.UnimplementedChatV1Server
}

// Create returns a fake chat id.
func (s *server) Create(_ context.Context, _ *srv.CreateRequest) (*srv.CreateResponse, error) {
	return &srv.CreateResponse{
		Id: utilites.Abs(gofakeit.Int64()),
	}, nil
}

// Delete does nothing and returns no error.
func (s *server) Delete(_ context.Context, _ *srv.DeleteRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// SendMessage returns a fake message id.
func (s *server) SendMessage(_ context.Context, _ *srv.SendMessageRequest) (*srv.SendMessageResponse, error) {
	return &srv.SendMessageResponse{
		Id: utilites.Abs(gofakeit.Int64()),
	}, nil
}

// main starts a grpc server that listens on port 50051 and implements the
// server-side rpc methods of the chat service.
func main() {
	addr, err := net.ResolveTCPAddr("tcp", host+grpcPort)
	if err != nil {
		log.Fatalf("failed to resolve address: %v", err)
	}

	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err = gofakeit.Seed(time.Now().UnixNano()); err != nil {
		log.Default().Printf("failed to seed random number generator: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	srv.RegisterChatV1Server(s, &server{})
	log.Printf("server listening at %v", listen.Addr())

	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
