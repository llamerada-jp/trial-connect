package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	v1 "github.com/llamerada-jp/trial-connect/proto/v1"
	"github.com/llamerada-jp/trial-connect/proto/v1/v1connect"
)

func main() {
	mux := http.NewServeMux()

	path, handler := v1connect.NewStreamServiceHandler(&streamImpl{})
	mux.Handle(path, handler)
	path, handler = v1connect.NewUnaryServiceHandler(&unaryImpl{})
	mux.Handle(path, handler)

	mux.Handle("/", http.FileServer(http.Dir("static")))

	go func() {
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			panic(err)
		}
	}()

	err := http.ListenAndServeTLS(":8443", "localhost.crt", "localhost.key", mux)
	if err != nil {
		panic(err)
	}
}

type streamImpl struct {
	v1connect.UnimplementedStreamServiceHandler
}

func (s *streamImpl) Echo(ctx context.Context, request *connect.Request[v1.StreamServiceEchoRequest], stream *connect.ServerStream[v1.StreamServiceEchoResponse]) error {
	message := request.Msg.GetMessage()

	for i := 0; i < 10; i++ {
		err := stream.Send(&v1.StreamServiceEchoResponse{
			Message: fmt.Sprintf("hello %s stream:%d", message, i),
		})
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

type unaryImpl struct {
	v1connect.UnimplementedUnaryServiceHandler
}

func (u *unaryImpl) Echo(ctx context.Context, request *connect.Request[v1.UnaryServiceEchoRequest]) (*connect.Response[v1.UnaryServiceEchoResponse], error) {
	fmt.Printf("method: %s\naddr: %s\nprotocol: %s\n", request.HTTPMethod(), request.Peer().Addr, request.Peer().Protocol)
	fmt.Println("header:")
	for k, vs := range request.Header() {
		fmt.Printf("  %s: %v\n", k, vs)
	}
	fmt.Println("query:")
	for k, vs := range request.Peer().Query {
		fmt.Printf("  %s: %v\n", k, vs)
	}

	message := request.Msg.GetMessage()
	return &connect.Response[v1.UnaryServiceEchoResponse]{
		Msg: &v1.UnaryServiceEchoResponse{
			Message: fmt.Sprintf("hello %s unary", message),
		},
	}, nil
}
