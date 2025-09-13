// Package main provides the entry point for the chat server.
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/based-chat/chat-server/utilites/mathx"
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	srv "github.com/based-chat/chat-server/pkg/chat/v1"
)

const (
	grpcPort = "50052"
	grpcHost = "localhost"

	errNotEnoughUsers = "not enough users"
	errInvalidID      = "invalid id"
	errEmptySender    = "empty sender"
	errEmptyMessage   = "empty message"
)

var (
	errFailedListen = errors.New("failed to listen")
	errFailedServe  = errors.New("failed to serve")
	errFailedSeed   = errors.New("failed to seed fakeit")
)

type server struct {
	srv.UnimplementedChatV1Server
}

// Create creates a new chat with the given usernames. The chat id is
// randomly generated and returned in the response. If the number of
// usernames is less than 2, then an InvalidArgument error is returned.
// The error is populated with the error message "not enough users".
func (s *server) Create(_ context.Context, req *srv.CreateRequest) (*srv.CreateResponse, error) {
	if len(req.GetUsernames()) < 2 {
		return nil, status.Error(codes.InvalidArgument, errNotEnoughUsers)
	}

	return &srv.CreateResponse{
		Id: mathx.Abs(gofakeit.Int64()),
	}, nil
}

// Delete deletes a chat with the given id. If the id is less than 0, then an
// InvalidArgument error is returned. The error is populated with the error
// message "invalid id".
func (s *server) Delete(_ context.Context, req *srv.DeleteRequest) (*emptypb.Empty, error) {
	if req.GetId() < 0 {
		return nil, status.Error(codes.InvalidArgument, errInvalidID)
	}

	return &emptypb.Empty{}, nil
}

// SendMessage returns a fake message id.
func (s *server) SendMessage(_ context.Context, req *srv.SendMessageRequest) (*srv.SendMessageResponse, error) {
	if req.GetChatId() < 0 {
		return nil, status.Error(codes.InvalidArgument, errInvalidID)
	}

	if len(req.GetSender()) == 0 {
		return nil, status.Error(codes.InvalidArgument, errEmptySender)
	}

	if len(req.GetMessage()) == 0 {
		return nil, status.Error(codes.InvalidArgument, errEmptyMessage)
	}

	return &srv.SendMessageResponse{
		Id: mathx.Abs(gofakeit.Int64()),
	}, nil
}

// main starts the grpc server and listens on the specified address.
// It seeds the random number generator and registers the user service.
// It then serves the grpc server and logs any errors that occur during serving.
func main() {
	// Listen on the specified address
	var lc net.ListenConfig

	listen, err := lc.Listen(context.Background(), "tcp", net.JoinHostPort(grpcHost, grpcPort))
	if err != nil {
		log.Fatalf("%s: %v", errFailedListen.Error(), err)
	}

	// Seed the random number generator
	if err := gofakeit.Seed(time.Now().UnixNano()); err != nil {
		log.Fatalf("%s: %v", errFailedSeed.Error(), err)
	}

	// Start the grpc server
	s := grpc.NewServer()
	reflection.Register(s)
	srv.RegisterChatV1Server(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatalf("%s: %v", errFailedServe.Error(), err)
	}
}
