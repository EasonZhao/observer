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
	reply.Interrupt = false
	fmt.Println(in.Timestamp)
	fmt.Println(in.Symbol)
	fmt.Println(in.Mining)
	fmt.Println(in.Confirm)
	fmt.Println(in.AddressTo)
	fmt.Println(in.Amount)
	fmt.Println(in.Txid)
	return &reply, nil
}

func main1() {
	lis, err := net.Listen("tcp", ":10889")
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
