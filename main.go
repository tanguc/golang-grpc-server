package main

import (
	"context"
	"fmt"
	pb "golang-grpc-server/upstream"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUpstreamPeerServiceServer
}

func (s *server) BidirectionalStreaming(server pb.UpstreamPeerService_BidirectionalStreamingServer) error {
	msg, err := server.Recv()
	fmt.Printf("before err : %v", err)
	for err == nil {
		fmt.Printf("Got input request = %v\n", msg)

		err = server.RecvMsg(msg)
	}
	fmt.Printf("Stream closed : %v", err)

	return nil
}

func (s *server) Ready(context.Context, *pb.ReadyRequest) (*pb.ReadyResult, error) {
	fmt.Printf("Ready?")
	result := new(pb.ReadyResult)
	header := generateHeader()

	result.Header = header
	result.Ready = true

	return result, nil
}

func generateHeader() *pb.Header {
	header := new(pb.Header)
	header.Address = "127.0.0.1"
	header.ClientUuid = "239-324-323-J392-32J4"
	header.Time = "08:49:11"

	return header
}

func (s *server) Live(server pb.UpstreamPeerService_LiveServer) error {
	fmt.Printf("Live?")
	fmt.Printf("yes I'm alive")

	msg, err := server.Recv()
	for ; err == nil; msg, err = server.Recv() {
		fmt.Printf("Got live request : %v", msg)

		header := generateHeader()
		server.Send(
			&pb.LiveResult{
				Header: header,
			},
		)
	}

	return err
}

func main() {
	fmt.Println("lol")

	address := "127.0.0.1:45888"

	lis, err := net.Listen("tcp", address)

	if err != nil {
		fmt.Printf("Error %v\n", err)
	}

	s := grpc.NewServer()
	pb.RegisterUpstreamPeerServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		fmt.Printf("Failed to serve : %v\n", err)
	}
	fmt.Printf("Finished")
}
