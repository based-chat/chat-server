// Package main provides the entry point for the chat server.
package main

import (
	"context"
	"log"
	"net"
	"time"

	srv "github.com/based-chat/chat-server/pkg/chat/v1"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	grpcPort = ":50051"
	host     = "0.0.0.0"
)

type server struct {
	srv.UnimplementedChatV1Server
}

func (s *server) Create(_ context.Context, _ *srv.CreateRequest) (*srv.CreateResponse, error) {

	gofakeit.Seed(time.Now().UnixNano())

	return &srv.CreateResponse{
		Id: abs(gofakeit.Int64()),
	}, nil
}

func (s *server) Delete(_ context.Context, _ *srv.DeleteRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(_ context.Context, _ *srv.SendMessageRequest) (*srv.SendMessageResponse, error) {

	return &srv.SendMessageResponse{
		Id: abs(gofakeit.Int64()),
	}, nil
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", host+grpcPort)
	if err != nil {
		log.Fatalf("failed to resolve address: %v", err)
	}

	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	srv.RegisterChatV1Server(s, &server{})
	log.Printf("server listening at %v", listen.Addr())

	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
