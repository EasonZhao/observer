package main

import (
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	pb "observer/protocol"
	"os"
	"time"
)

const (
	address     = "localhost:10087"
	defaultName = "listaddrs"
)

func newApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "observer client"
	app.Name = "obs-cli"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "listaddresses",
			Action: listAddresses,
			Usage:  "list mananger address",
		},
	}
	return app
}

func listAddresses(c *cli.Context) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	processor := pb.NewProcessorClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := processor.ListAddrs(ctx, &pb.ListAddrsRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for _, v := range r.Addrs {
		log.Fatalf(v)
	}
	return nil
}

func main() {
	app := newApp()
	app.Run(os.Args)
}
