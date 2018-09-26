package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	pb "observer/protocol"
)

// Processor comment
type Processor struct {
	addrMgr AddressManager
	server  *grpc.Server
}

// AddressManager interface for manager address
type AddressManager interface {
	Bind() string
	Addresses() map[string]string
}

// NewProcessor create new application
func NewProcessor(mgr AddressManager) *Processor {
	processor := &Processor{addrMgr: mgr}
	return processor
}

// ListAddrs comment
func (s *Processor) ListAddrs(ctx context.Context, in *pb.ListAddrsRequest) (*pb.ListAddrsReply, error) {
	reply := pb.ListAddrsReply{}
	addrs := s.addrMgr.Addresses()
	reply.Addrs = make([]string, len(addrs))
	i := 0
	for k := range addrs {
		reply.Addrs[i] = k
		i++
	}
	return &reply, nil
}

// AsyncGRPC asynchronous start grpc
func (s *Processor) AsyncGRPC(c *chan error) error {
	//start grpc
	lis, err := net.Listen("tcp", s.addrMgr.Bind())
	if err != nil {
		return err
	}
	s.server = grpc.NewServer()
	pb.RegisterProcessorServer(s.server, s)
	// register
	reflection.Register(s.server)
	go func() {
		if err := s.server.Serve(lis); err != nil {
			*c <- err
		}
		*c <- nil
	}()
	return nil
}

// StopGRPC stop grpc service
func (s *Processor) StopGRPC() {
	if s.server != nil {
		s.server.Stop()
		s.server = nil
	}
}
