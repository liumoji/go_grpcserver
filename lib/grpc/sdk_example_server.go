package grpc

import (
	"context"
	"fmt"
	"io"
	"sync"

	pb "SDK/proto/grpcexample"
)

type ExampleServer struct {
	pb.UnimplementedExampleServiceServer

	mu          sync.Mutex // protects routeNotes
	streamNotes map[int]pb.ExampleService_EchoStreamServer
	uuid        int
}

func (s *ExampleServer) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	var testP *pb.EchoResponse
	testP = &pb.EchoResponse{}
	testP.Message = "this sdk send "
	testP.Code = 1
	return testP, nil
}

func (s *ExampleServer) EchoStream(stream pb.ExampleService_EchoStreamServer) error {
	var thisId int
	s.mu.Lock()
	s.uuid++
	thisId = s.uuid
	s.streamNotes[thisId] = stream
	s.mu.Unlock()
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Print("EchoStream io.EOF uuid: " + fmt.Sprint(thisId))
			delete(s.streamNotes, thisId)
			return nil
		}
		if err != nil {
			log.Printf("EchoStream err : %s , uuid: %d end", err, thisId)
			delete(s.streamNotes, thisId)
			return err
		}

		log.Info("Receive client message: ID: " + fmt.Sprint(in.ID) + " Message: " + in.Message)
		sendMessage := &pb.EchoResponse{Message: "server send from uuid: " + fmt.Sprint(thisId) + " message: this is server!"}
		serr := stream.Send(sendMessage)
		if serr != nil {
			log.Printf("Send end")
			delete(s.streamNotes, thisId)
			return serr
		}
	}
}

func newServer() *ExampleServer {
	s := &ExampleServer{streamNotes: make(map[int]pb.ExampleService_EchoStreamServer)}
	return s
}
