# grpc-go-messages
A simple grpc client and server example in Go

TODO: Update this readme!

`go get github.com/Luke-Davies/grpc-go-messages`

Server:

`cd $GOPATH/src/github.com/Luke-Davies/grpc-go-messages`

`go run server/server.go`

Client:

`cd $GOPATH/src/github.com/Luke-Davies/grpc-go-messages`

`go run client/client.go AddMessage 'test message'`

`go run client/client.go GetMessages`

`go run client/client.go GetMessages 1000`

`go run client/client.go DeleteMessages`
