// Package main provides the entry point for the chat client.
package main

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	srv "github.com/based-chat/chat-server/pkg/chat/v1"
)

const (
	grpcPort       = ":50051"
	host           = "0.0.0.0"
	maxTimemeout   = 3 * time.Second
	fakeCountUsers = 5
)

func main() {
	conn, err := grpc.NewClient(host+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Default().Println("failed to connect: " + err.Error())
	}
	defer func() { _ = conn.Close() }() // nolintconn.Close()

	c := srv.NewChatV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), maxTimemeout)
	defer cancel()

	usernames := make([]string, fakeCountUsers)
	for i := range usernames {
		usernames[i] = gofakeit.Username()
	}

	_, err = c.Create(ctx, &srv.CreateRequest{Usernames: usernames})
	if err != nil {
		log.Default().Println("failed to create: " + err.Error())
	}

	_, err = c.SendMessage(ctx, &srv.SendMessageRequest{
		ChatId:  abs(gofakeit.Int64()),
		Sender:  gofakeit.Username(),
		Message: gofakeit.Phrase(),
	})
	if err != nil {
		log.Default().Println("failed to send message: " + err.Error())
	}

	_, err = c.Delete(ctx, &srv.DeleteRequest{
		Id: abs(gofakeit.Int64()),
	})
	if err != nil {
		log.Default().Println("failed to delete: " + err.Error())
	}
}

func abs(x int64) int64 {
	if x == math.MinInt64 {
		return math.MaxInt64
	}
	if x < 0 {
		return -x
	}
	return x
}
