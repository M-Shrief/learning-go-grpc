package main

import (
	"context"
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

func main() {
	con, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Couldn't connect to: %v", err)
	}
	defer con.Close()

	client := pb.NewPingClient(con)

	clientPing(client)
}
