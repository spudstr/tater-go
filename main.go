package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	stop := time.NewTicker(70 * time.Second)

	for {
		select {
		case <-stop.C:
			err := stream.CloseSend()
			if err != nil {
				log.Fatalf("can not close stream %v", err.Error())
			}
			return
		case <-stream.Context().Done():
			log.Fatalf("can not stream data %v", stream.Context().Err())
		default:
			data, err := stream.Recv()
			if err != nil {
				panic(err)
			}
			fmt.Printf("%+v\n", data)

		}
	}

}
