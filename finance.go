package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "observer/finance"
	"time"
)

func depositNotify(req pb.DepositRequest, host string) (*pb.DepositReply, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewFinanceServiceClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.DepositNotify(ctx, &req)
	if err != nil {
		return nil, err
	}
	return r, nil
}
