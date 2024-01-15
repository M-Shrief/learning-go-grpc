package main

import (
	"context"
	"io"
	"learning-go/pb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":8080"
)

func clientPing(client pb.PingClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.PingPong(ctx, &pb.PingRequest{Message: "Ping"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("%s", res.Message)
}

func clientChat(client pb.ChatClient) {
	log.Printf("Bidirectional Streaming started")
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("Couldn't send names: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while streaming %v", err)
			}
			log.Println(message)
		}
		close(waitc)
	}()

	// Greeting some people
	for _, name := range []string{"Satya", "Sumi", "Arya"} {
		req := &pb.ChatRequest{
			Name:    name,
			Message: "Hey!",
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	stream.CloseSend()
	<-waitc
	log.Printf("Bidirectional Streaming finished")
}

func main() {
	con, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Couldn't connect to: %v", err)
	}
	defer con.Close()

	client := pb.NewPingClient(con)

	clientPing(client)

	client2 := pb.NewChatClient(con)
	clientChat(client2)
}
