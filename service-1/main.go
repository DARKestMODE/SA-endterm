package main

import (
	"awesomeProject/service-1/proto/primepb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type Server struct {
	primepb.UnimplementedSumServiceServer
}

func (s *Server) PrimeDecomposition(req *primepb.PrimeDecompositionRequest, stream primepb.SumService_PrimeDecompositionServer) error{
	number := req.GetNumber().GetNumber()
	k := int32(2)
	for number > 1 {
		if number % k == 0 {
			res := &primepb.PrimeDecompositionResponse{Result: k}
			number = number/k
			if err := stream.Send(res); err != nil {
				log.Fatalf("error while sending PrimeDecomposition requests: %v", err.Error())
			}
			time.Sleep(time.Second)
		} else {
			k = k + 1
		}
	}
	return nil
}

func (s *Server) ComputeAverage(stream primepb.SumService_ComputeAverageServer) error {
	var result int32 = 0
	var i float64 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&primepb.ComputeAverageResponse{
				Result: float64(result)/i,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		i += 1
		result += req.Number.GetNumber()
	}
}


func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}

	s := grpc.NewServer()
	primepb.RegisterSumServiceServer(s, &Server{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}
