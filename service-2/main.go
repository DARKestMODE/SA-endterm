package main

import (
	"awesomeProject/service-2/proto/primepb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func doPrimeDecomposition(c primepb.SumServiceClient) {
	ctx := context.Background()
	request := &primepb.PrimeDecompositionRequest{Number: &primepb.Number{
		Number: 120,
	}}

	stream, err := c.PrimeDecomposition(ctx, request)

	if err != nil {
		log.Fatalf("error while calling PrimeDecomposition RPC %v", err)
	}
	defer stream.CloseSend()

	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("error while reciving from PrimeDecomposition RPC %v", err)
		}
		log.Printf("response from PrimeDecomposition RPC:%v \n", res.GetResult())
	}
}

func doComputeAverage(c primepb.SumServiceClient) {

	requests := []*primepb.ComputeAverageRequest{
		{
			Number: &primepb.Number{
				Number: 1,
			},
		},
		{
			Number: &primepb.Number{
				Number: 2,
			},
		},
		{
			Number: &primepb.Number{
				Number: 3,
			},
		},
		{
			Number: &primepb.Number{
				Number: 4,
			},
		},
	}

	ctx := context.Background()
	stream, err := c.ComputeAverage(ctx)
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res)
}


func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("micro-1:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := primepb.NewSumServiceClient(conn)

	doPrimeDecomposition(c)
	doComputeAverage(c)
}
