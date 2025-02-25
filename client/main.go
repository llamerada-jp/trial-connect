package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/llamerada-jp/trial-connect/proto/v1"
	"github.com/llamerada-jp/trial-connect/proto/v1/v1connect"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpClient := http.DefaultClient
	unaryClient := v1connect.NewUnaryServiceClient(httpClient, "http://localhost:8080")
	unaryRes, err := unaryClient.Echo(ctx, connect.NewRequest(&v1.UnaryServiceEchoRequest{
		Message: "world",
	}))
	if err != nil {
		panic(err)
	}
	fmt.Println(unaryRes.Msg.GetMessage())

	streamClient := v1connect.NewStreamServiceClient(httpClient, "http://localhost:8080")
	streamRes, err := streamClient.Echo(ctx, connect.NewRequest(&v1.StreamServiceEchoRequest{
		Message: "world",
	}))
	if err != nil {
		panic(err)
	}
	for streamRes.Receive() {
		res := streamRes.Msg()
		if err != nil {
			break
		}
		fmt.Println(res.GetMessage())
	}
}
