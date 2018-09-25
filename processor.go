package main

import (
	"golang.org/x/net/context"
	pb "observer/protocol"
)

// Processor comment
type Processor struct {
	app *Application
}

// NewProcessor create new application
func NewProcessor(app *Application) *Processor {
	return &Processor{app}
}

// ListAddrs comment
func (s *Processor) ListAddrs(ctx context.Context, in *pb.ListAddrsRequest) (*pb.ListAddrsReply, error) {
	reply := pb.ListAddrsReply{}
	reply.Addrs = make([]string, len(s.app.addresses))
	i := 0
	for k, _ := range s.app.addresses {
		reply.Addrs[i] = k
		i++
	}
	return &reply, nil
}
