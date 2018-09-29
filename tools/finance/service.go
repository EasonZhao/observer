package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	pb "observer/finance"
)

// Service comment
type Service struct {
}

// DepositNotify comment
func (s *Service) DepositNotify(ctx context.Context, in *pb.DepositRequest) (*pb.DepositReply, error) {
	reply := pb.DepositReply{}
	if in.Confirm >= 6 {
		reply.Interrupt = true
	} else {
		fmt.Println(in.Confirm)
	}
	return &reply, nil
}

func main() {
	lis, err := net.Listen("tcp", ":10089")
	if err != nil {
		log.Fatalln(err)
	}
	server := grpc.NewServer()
	s := &Service{}
	pb.RegisterFinanceServiceServer(server, s)
	// register
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
