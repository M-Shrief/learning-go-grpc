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

func ping(client pb.PingClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.PingPong(ctx, &pb.PingRequest{Message: "Ping"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("%s", res.Message)
}

func computeAverage(client pb.CalculatorClient) {
	log.Println("Starting to do a ComputeAverage Client Streaming RPC...")

	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	for _, number := range []int32{3, 5, 7, 9, 54, 76} {
		log.Printf("Sending number: %v\n", number)
		stream.Send(&pb.ComputeAverageRequest{Number: number})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}
	log.Printf("The Average is: %v\n", res.GetAverage())
}

func primeNumberDecomposition(client pb.CalculatorClient) {
	log.Println("Starting to do a PrimeDecomposition Server Streaming RPC...")
	req := &pb.PrimeNumberDecompositionRequest{Number: 12390392840}

	stream, err := client.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeDecomposition RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		log.Println("Prime Factor: %v", res.GetPrimeFactor())
	}
}

func chat(client pb.ChatClient) {
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

	// client := pb.NewPingClient(con)

	// ping(client)

	// client2 := pb.NewChatClient(con)
	// chat(client2)

	client3 := pb.NewCalculatorClient(con)
	// computeAverage(client3)
	primeNumberDecomposition(client3)
}
