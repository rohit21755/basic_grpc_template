package main

import (
	"context"
	"errors"
	"fmt"
	pb "grpc-calculator/calculatorpb"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCalculatorServiceServer
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	result := req.A + req.B
	return &pb.AddResponse{Result: result}, nil
}

func (s *server) Subtract(ctx context.Context, req *pb.SubtractRequest) (*pb.SubtractResponse, error) {
	if req.A > req.B {
		return &pb.SubtractResponse{Result: req.A - req.B}, nil
	} else {
		return nil, errors.New("Invalid value")
	}
}

func (s *server) Multiply(ctx context.Context, req *pb.MultiplyRequest) (*pb.MultiplyResponse, error) {
	return &pb.MultiplyResponse{Result: req.A * req.B}, nil
}

func (s *server) Divide(ctx context.Context, req *pb.DivideRequest) (*pb.DivideResponse, error) {
	if req.B > 0 {
		return &pb.DivideResponse{Result: req.A / req.B}, nil
	} else {
		return nil, errors.New("Invalid values")
	}
}

// server streaming
func (s *server) PrimeFactors(req *pb.PrimeRequest, stream pb.CalculatorService_PrimeFactorsServer) error {
	n := req.Number
	divisor := int32(2)
	for n > 1 {
		if n%divisor == 0 {
			stream.Send(&pb.PrimeResponse{PrimeFactor: divisor})
			n = n / divisor
		} else {
			divisor++
		}
	}
	return nil
}

func (s *server) Average(stream pb.CalculatorService_AverageServer) error {
	var sum, count int32
	for {
		in, err := stream.Recv()
		fmt.Println(in)
		if err == io.EOF {
			avg := float64(sum) / float64(count)
			return stream.SendAndClose(&pb.AverageResponse{Average: avg})
		}
		if err != nil {
			return err
		}
		sum += in.Number
		count++
	}
}

func main() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("Failed to listen on :8000")
	}

	gerpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(gerpcServer, &server{})

	if err := gerpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve: %v", err)
	}
}
