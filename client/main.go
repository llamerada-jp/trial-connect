package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"connectrpc.com/connect"
	v1 "github.com/llamerada-jp/trial-connect/proto/v1"
	"github.com/llamerada-jp/trial-connect/proto/v1/v1connect"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url := "http://localhost:8080"
	if runtime.GOOS == "js" {
		url = "/"
	}

	httpClient := http.DefaultClient
	unaryClient := v1connect.NewUnaryServiceClient(httpClient, url)
	unaryRes, err := unaryClient.Echo(ctx, connect.NewRequest(&v1.UnaryServiceEchoRequest{
		Message: "world",
	}))
	if err != nil {
		panic(err)
	}
	fmt.Println(unaryRes.Msg.GetMessage())

	streamClient := v1connect.NewStreamServiceClient(httpClient, url)
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
