package main

import (
	"context"
	"fmt"
	pb "grpc-calculator/calculatorpb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conn.GetState().String())
	defer conn.Close()
	clinet := pb.NewCalculatorServiceClient(conn)
	// addResp, err := clinet.Add(context.Background(), &pb.AddRequest{A: 10, B: 20})
	// if err != nil {
	// 	log.Fatal("Error in adding two numbers", err)
	// }
	// fmt.Println("Add Result: ", addResp.Result)
	// fmt.Println("time took: ", time.Since(now))

	// server streaming
	time.Sleep(time.Second * 3)
	// stream, err := clinet.PrimeFactors(context.Background(), &pb.PrimeRequest{Number: 48})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for {
	// 	res, err := stream.Recv()
	// 	if err == io.EOF {
	// 		log.Fatal("Stream error: ", err)
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("Stream error: %v", err) // Actual errors
	// 	}
	// 	time.Sleep(time.Second)
	// 	fmt.Print(res.PrimeFactor, "  \n")
	// }
	// fmt.Println()
	//client streaming
	cStream, _ := clinet.Average(context.Background())
	for _, val := range []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		time.Sleep(time.Second * 2)
		cStream.Send(&pb.Number{Number: val})
	}
	avgResp, _ := cStream.CloseAndRecv()
	fmt.Println("Average: ", avgResp.Average)
}
