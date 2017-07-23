package main

import (
	"log"
	"net"
	"os"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"errors"

	"fmt"

	pb "github.com/luke-davies/grpc-go-messages/messages"
)

// MessagesService definition
type MessagesService struct {
	storedMessages []*pb.Message
	nextID         int64
}

// NewMessagesService return new instance of server
func NewMessagesService() *MessagesService {
	return &MessagesService{nextID: 1000}
}

func (s *MessagesService) genNextID() int64 {
	defer func() { s.nextID++ }()
	return s.nextID
}

// GetMessages returns all messages
func (s *MessagesService) GetMessages(context.Context, *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	return &pb.GetMessagesResponse{Messages: s.storedMessages}, nil
}

//AddMessage adds a new message
func (s *MessagesService) AddMessage(ctx context.Context, req *pb.AddMessageRequest) (*pb.AddMessageResponse, error) {
	newMsg := pb.Message{
		Id:   s.genNextID(),
		Text: req.GetText(),
	}
	s.storedMessages = append(s.storedMessages, &newMsg)
	resp := &pb.AddMessageResponse{Id: newMsg.GetId()}
	return resp, nil
}

// GetMessage get a single message by Id
func (s *MessagesService) GetMessage(ctx context.Context, req *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	for _, v := range s.storedMessages {
		if v.GetId() == req.GetId() {
			return &pb.GetMessageResponse{Text: v.GetText()}, nil
		}
	}
	errMsg := fmt.Sprintf("No message found for ID %v", req.GetId())
	return &pb.GetMessageResponse{Text: errMsg}, errors.New(errMsg)
}

// DeleteMessages deletes all messages
func (s *MessagesService) DeleteMessages(context.Context, *pb.DeleteMessagesRequest) (*pb.DeleteMessagesRespone, error) {
	s.storedMessages = []*pb.Message{}
	return &pb.DeleteMessagesRespone{}, nil
}

func main() {
	port := ":3000"
	if v := os.Getenv("PORT"); v != "" {
		port = ":" + v
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMessagesServiceServer(s, NewMessagesService())
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
