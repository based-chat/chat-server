// Package main provides the entry point for the chat client.
package main

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/based-chat/chat-server/utilites"

	srv "github.com/based-chat/chat-server/pkg/chat/v1"
)

const (
	grpcPort       = ":50051"
	host           = "0.0.0.0"
	maxTimeout     = 3 * time.Second
	fakeCountUsers = 5
)

// main demonstrates basic usage of the chat service rpc api.
// It connects to the server, creates a chat with a few users,
// sends a message, and then deletes the chat.
func main() {
	conn, err := grpc.NewClient(host+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer func() {
		connErr := conn.Close()
		if connErr != nil {
			log.Default().Println("failed to close connection: " + connErr.Error())
		}
	}()

	c := srv.NewChatV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	usernames := make([]string, fakeCountUsers)
	for i := range usernames {
		usernames[i] = gofakeit.Username()
	}

	createResponse, err := c.Create(ctx, &srv.CreateRequest{Usernames: usernames})
	if err != nil {
		log.Default().Println("failed to create: " + err.Error())
		return
	}

	_, err = c.SendMessage(ctx, &srv.SendMessageRequest{
		ChatId:  createResponse.GetId(),
		Sender:  gofakeit.Username(),
		Message: gofakeit.Phrase(),
	})
	if err != nil {
		log.Default().Println("failed to send message: " + err.Error())
		return
	}

	_, err = c.Delete(ctx, &srv.DeleteRequest{
		Id: utilites.Abs(createResponse.GetId()),
	})
	if err != nil {
		log.Default().Println("failed to delete: " + err.Error())
		return
	}
}
