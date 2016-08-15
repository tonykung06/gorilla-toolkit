package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

func main() {
	go startServer()
	go startClient()
	fmt.Scanln()
}

func startClient() {
	args := ClientArgs{4, 5}
	msg, _ := json.EncodeClientRequest("ArithService.Add", args)
	resp, _ := http.Post("http://localhost:1234/rpc", "application/json", bytes.NewReader(msg))
	var result int
	json.DecodeClientResponse(resp.Body, &result)
	fmt.Println(result)
}

type ClientArgs struct {
	A, B int
}

func startServer() {
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterService(new(ArithService), "")

	server.RegisterBeforeFunc(func(ri *rpc.RequestInfo) {
		fmt.Println("before", ri)
	})

	server.RegisterAfterFunc(func(ri *rpc.RequestInfo) {
		fmt.Println("after", ri)
	})

	http.Handle("/rpc", server)
	httpServer := &http.Server{
		Addr:    ":1234",
		Handler: nil,
		ConnState: func(conn net.Conn, state http.ConnState) {
			// fmt.Println(state)
		},
	}
	httpServer.ListenAndServe()
}

type Args struct {
	A, B int
}

type ArithService struct{}

func (a *ArithService) Add(req *http.Request, args *Args, reply *int) error {
	*reply = args.A + args.B
	fmt.Println("Calculating...")
	return nil
}
