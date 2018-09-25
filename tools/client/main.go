package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	pb "observer/protocol"
	"time"
)

const (
	address     = "localhost:10087"
	defaultName = "listaddrs"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProcessorClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ListAddrs(ctx, &pb.ListAddrsRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for _, v := range r.Addrs {
		fmt.Println(v)
	}
}
