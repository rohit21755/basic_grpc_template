package main

import (
	"context"
	"fmt"
	pb "grpc-calculator/calculatorpb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	clinet := pb.NewCalculatorServiceClient(conn)
	addResp, err := clinet.Add(context.Background(), &pb.AddRequest{A: 10, B: 20})
	if err != nil {
		log.Fatal("Error in adding two numbers", err)
	}
	fmt.Println("Add Result: ", addResp.Result)
}
