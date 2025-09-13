// Package main provides the entry point for the chat client.
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
	"google.golang.org/grpc/credentials/insecure"

	srv "github.com/based-chat/chat-server/pkg/chat/v1"
)

const (
	grpcPort       = "50052"
	grpcHost       = "localhost"
	maxTimeout     = 1 * time.Second
	fakeCountUsers = 5
)

var (
	errFailedConnect         = errors.New("failed to connect: %v")
	errFailedCloseConnection = "failed to close connection: %v"
	errFailedCreateChat      = "failed to create chat: %v"
	errFailedSendMessage     = "failed to send message: %v"
	errFailedDeleteChat      = "failed to delete chat: %v"
)

// main demonstrates basic usage of the chat service rpc api.
// It connects to the server, creates a chat with a few users,
// sends a message, and then deletes the chat.
func main() {
	// Connect to the grpc server
	addr := net.JoinHostPort(grpcHost, grpcPort)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Default().Printf(errFailedConnect.Error(), err)
		return
	}

	defer func() {
		connErr := conn.Close()
		if connErr != nil {
			log.Default().Printf(errFailedCloseConnection, connErr)
		}
	}()

	c := srv.NewChatV1Client(conn)

	usernames := make([]string, fakeCountUsers)
	for i := range usernames {
		usernames[i] = gofakeit.Username()
	}

	ctxCreateResponse, cancelCreateResponse := context.WithTimeout(context.Background(), maxTimeout)

	createResponse, err := c.Create(ctxCreateResponse, &srv.CreateRequest{Usernames: usernames})
	if err != nil {
		log.Default().Printf(errFailedCreateChat, err)
		return
	}

	cancelCreateResponse()

	ctxSendMessage, cancelSendMessage := context.WithTimeout(context.Background(), maxTimeout)

	_, err = c.SendMessage(ctxSendMessage, &srv.SendMessageRequest{
		ChatId:  createResponse.GetId(),
		Sender:  gofakeit.Username(),
		Message: gofakeit.Phrase(),
	})
	if err != nil {
		log.Default().Printf(errFailedSendMessage, err)
		return
	}

	cancelSendMessage()

	ctxDelete, cancelDelete := context.WithTimeout(context.Background(), maxTimeout)

	_, err = c.Delete(ctxDelete, &srv.DeleteRequest{
		Id: mathx.Abs(createResponse.GetId()),
	})
	if err != nil {
		log.Default().Printf(errFailedDeleteChat, err)
		return
	}

	cancelDelete()
}
