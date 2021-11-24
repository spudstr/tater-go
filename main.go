package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/spudstr/tater-go/schema"
	"google.golang.org/grpc"
)

func main() {
	domain := os.Getenv("quodd_domain")
	username := os.Getenv("quodd_username")
	password := os.Getenv("quodd_password")
	gateway := os.Getenv("quodd_gateway")
	server := "gateway.quodd.com:55012"

	if domain == "" || username == "" || password == "" || gateway == "" {
		panic("Missing environment variables")
	}
	// open grpc connection to quodd
	// dial server
	conn, err := grpc.Dial(server, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}
	defer conn.Close()

	// populate Authrequest
	var authReq pb.AuthRequest
	authReq.Username = []byte(username)
	authReq.Password = []byte(password)
	authReq.Domain = []byte(domain)

	stub := pb.NewTableApiClient(conn)

	// call auth
	authResp, err := stub.Authenticate(context.Background(), &authReq)
	if err != nil {
		log.Fatalf("can not auth with server %v", err)
	}
	fmt.Printf("%+v\n", authResp)

	// call get
	var streamRequest pb.StreamRequest
	var keys [][]byte
	//	keys = append(keys, []byte("O:TSLA*"))
	keys = append(keys, []byte("O:SPY*"))

	streamRequest.SessionToken = authResp.SessionToken
	streamRequest.Keys = keys

	stream, _ := stub.StreamData(context.Background(), &streamRequest)
	chanMessages := make(chan interface{}, 10000)

	go func() {
		for msgBytes := range chanMessages {
			fmt.Printf("%+v\n", msgBytes)
			fmt.Printf("%d\n", len(chanMessages))
		}
	}()

	for {
		var msg interface{}
		msg, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("can not read stream %v", err)
		}
		chanMessages <- msg
	}

}
