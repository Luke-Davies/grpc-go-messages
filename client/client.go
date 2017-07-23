package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"strconv"

	pb "github.com/luke-davies/grpc-go-messages/messages"
	"google.golang.org/grpc"
)

// AvailableCommands - the commands this service expects
const (
	defaultPort       = ":3000"
	AvailableCommands = "GetMessages, AddMessage, GetMessage, DeleteMessages"
)

// Close - close the given io
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func callGetMessages(c pb.MessagesServiceClient) {
	r, err := c.GetMessages(context.Background(), &pb.GetMessagesRequest{})
	if err != nil {
		log.Fatalf("could not get messages: %v", err)
	}
	fmt.Println(r.GetMessages())
}

func callAddMessage(c pb.MessagesServiceClient, messageText string) {
	r, err := c.AddMessage(context.Background(), &pb.AddMessageRequest{Text: messageText})
	if err != nil {
		log.Fatalf("could not add message: %v", err)
	}
	fmt.Println(r.GetId())
}

func callGetMessage(c pb.MessagesServiceClient, msgID int64) {
	r, err := c.GetMessage(context.Background(), &pb.GetMessageRequest{Id: msgID})
	if err != nil {
		log.Fatalf("could not get message: %v", err)
	}
	fmt.Println(r.GetText())
}

func callDeleteMessages(c pb.MessagesServiceClient) {
	_, err := c.DeleteMessages(context.Background(), &pb.DeleteMessagesRequest{})
	if err != nil {
		log.Fatalf("could not delete messages: %v", err)
	}
	fmt.Println("Messages Deleted")
}

func getAddress() string {
	port := defaultPort
	if v := os.Getenv("PORT"); v != "" {
		port = ":" + v
	}
	address := fmt.Sprintf("localhost%v", port)
	return address
}

func main() {

	address := getAddress()

	// Set up a connection to the server.
	conn, dialErr := grpc.Dial(address, grpc.WithInsecure())
	if dialErr != nil {
		log.Fatalf("did not connect: %v", dialErr)
	}
	defer Close(conn)
	c := pb.NewMessagesServiceClient(conn)

	if len(os.Args) == 1 {
		log.Fatalln("No command provided. Please provide one of ", AvailableCommands)
	}
	switch command := os.Args[1]; command {
	case "GetMessages":
		callGetMessages(c)
	case "AddMessage":
		messageText := ""
		if len(os.Args) > 2 {
			messageText = os.Args[2]
		}
		callAddMessage(c, messageText)
	case "GetMessage":
		var msgID int64
		if len(os.Args) > 2 {
			var err error
			msgID, err = strconv.ParseInt(os.Args[2], 10, 64)
			if err != nil {
				log.Fatalf("could not parse message id (expecting int64): %v", err)
			}
		} else {
			log.Fatalln("Insufficient number of arguments. No message id given")
		}
		callGetMessage(c, msgID)
	case "DeleteMessages":
		callDeleteMessages(c)
	default:
		log.Fatalln("Unrecognsed command. Expected one of ", AvailableCommands)
	}
}
