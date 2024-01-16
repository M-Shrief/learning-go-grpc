package main

import (
	"context"
	"io"
	"learning-go/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

type server struct {
	pb.UnimplementedPingServer
	pb.UnimplementedChatServer
	pb.UnimplementedCalculatorServer
}

func (s *server) PingPong(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	log.Println(req.GetMessage())
	return &pb.PongResponse{Message: "Pong"}, nil
}

func (s *server) ComputeAverage(stream pb.Calculator_ComputeAverageServer) error {
	log.Println("Starting Computing the Average...")

	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			log.Printf("Average is: %v", average)
			return stream.SendAndClose(&pb.ComputeAverageResponse{Average: average})

		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		log.Printf("Receiving: %v", req.GetNumber())
		sum += req.GetNumber()
		count++
	}
}

func (s *server) Chat(stream pb.Chat_ChatServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
		log.Printf("Request : %v", req.Name)

		res := &pb.ChatResponse{
			Message: "Hello " + req.Name,
		}

		if err := stream.Send(res); err != nil {
			return err
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	grpcServer := grpc.NewServer()
	// pb.RegisterPingServer(grpcServer, &server{})
	// pb.RegisterChatServer(grpcServer, &server{})
	pb.RegisterCalculatorServer(grpcServer, &server{})
	log.Printf("Server starter at: %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

}
