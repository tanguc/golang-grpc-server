package main

import (
	"context"
	"errors"
	"fmt"
	pb "golang-grpc-server/proto"
	"log"
	"net"
	"os"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/urfave/cli/v2"
)

type server struct {
	pb.UnimplementedUpstreamPeerServiceServer
	ready_sleep_ms int64
}

func (s *server) BidirectionalStreaming(server pb.UpstreamPeerService_BidirectionalStreamingServer) error {
	// var msg, err
	var msg *pb.InputStreamRequest
	var err error
	fmt.Println("bi-directional streaming ON")

	for ; err == nil; msg, err = server.Recv() {
		if msg != nil {
			fmt.Printf("Streaming INPUT = %v\n", msg)

			payload := []byte("Salut à toi cher downstream")
			// time.Sleep(1000 * time.Millisecond)
			result := pb.OutputStreamRequest{
				Header:  generateHeader(msg.GetHeader()),
				Payload: payload,
			}
			server.Send(&result)
		}
	}
	fmt.Printf("Stream closed : %v\n", err)

	return nil
}

func (s *server) Ready(context context.Context, in *emptypb.Empty) (*pb.ReadyResult, error) {
	fmt.Printf("Ready?")

	if s.ready_sleep_ms > 0 {
		time.Sleep(time.Duration(s.ready_sleep_ms) * time.Millisecond)
	}

	result := new(pb.ReadyResult)

	result.Time = time.Now().Local().Format("Mon Jan 2 15:04:05 MST 2006")
	result.Ready = true

	return result, nil
}

func generateHeader(oldHeader *pb.Header) *pb.Header {
	header := new(pb.Header)
	if oldHeader != nil {
		header.Address = oldHeader.Address
		header.ClientUuid = oldHeader.ClientUuid
		header.Time = time.Now().Local().Format("Mon Jan 2 15:04:05 MST 2006")
	} else {
		fmt.Printf("Warning: the old header was nil")

	}

	return header
}

func (s *server) Live(server pb.UpstreamPeerService_LiveServer) error {
	fmt.Printf("Live?\n")
	fmt.Printf("yes I'm alive\n")

	// fmt.Printf("Context value of live : %v\n", s.context)

	msg, err := server.Recv()
	for ; err == nil; msg, err = server.Recv() {
		fmt.Printf("Got live request : %v", msg)

		server.Send(
			&pb.LiveResult{
				Time: time.Now().Local().Format("Mon Jan 2 15:04:05 MST 2006"),
				Live: true,
			},
		)
	}

	return err
}

func main() {

	flagPortName := "port"
	var flagPortValue int64

	flagReadySleep := "ready_sleep"
	var flagReadySleepMs int64

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.Int64Flag{
			Name:     flagPortName,
			Value:    45888,
			Usage:    "Server port to listen",
			Required: true,
		},
		&cli.Int64Flag{
			Name:     flagReadySleep,
			Value:    0,
			Usage:    "Number of milliseconds to wait before to send readiness",
			Required: false,
		},
	}

	app.Usage = "Start a GRPC server with predefined proto"

	app.Action = func(c *cli.Context) error {
		flagPortValue = c.Value(flagPortName).(int64)

		if !(flagPortValue > 1024 && flagPortValue <= 65535) {
			return errors.New("Port should be between 1024 and 65535")
		}

		flagReadySleepMs = c.Value(flagReadySleep).(int64)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	address := fmt.Sprintf("127.0.0.1:%v", flagPortValue)
	fmt.Printf("Starting server on [%v]\n", address)

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Failed to listen for server, cause: %v\n", err)
	}

	server := new(server)
	server.ready_sleep_ms = flagReadySleepMs

	s := grpc.NewServer()
	pb.RegisterUpstreamPeerServiceServer(s, server)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve GRPC, cause: %v\n", err)
	}
	fmt.Println("Server terminated")
}
