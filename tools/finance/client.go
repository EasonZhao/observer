package main

import (
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	pb "observer/finance"
	"os"
	"time"
)

const (
	address = "localhost:10889"
)

func newApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "finace client"
	app.Name = "obs-cli"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "depositnotify",
			Action: depositNotify,
			Usage:  "deposit notify",
		},
	}
	return app
}

func depositNotify(c *cli.Context) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewFinanceServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.DepositRequest{
		Timestamp: time.Now().Unix(),
		Symbol:    "btc",
		Amount:    "0.00001",
		Txid:      "0x111122334444",
		Mining:    false,
		Confirm:   1,
		AddressTo: "0x312121323434234234254dsfg22344",
	}
	r, err := client.DepositNotify(ctx, req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(r.Interrupt)
	return nil
}

func main1() {
	app := newApp()
	app.Run(os.Args)
}
