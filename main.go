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

type grc_details struct {
	domain   string
	username string
	password string
}

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
	defer conn.Close()

	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

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
	keys = append(keys, []byte("O:TSLA*"))
	keys = append(keys, []byte("O:SPY*"))

	streamRequest.SessionToken = authResp.SessionToken
	streamRequest.Keys = keys

	stream, err := stub.StreamData(context.Background(), &streamRequest)
	if err != nil {
		log.Fatalf("can not get with server %v", err)
	}

	done := make(chan bool)
	counter := int64(0)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true //close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			fmt.Printf("%+v\n", counter)
			counter++
			log.Printf("1")
			_ = resp
		}
	}()

	<-done
	log.Printf("Finished")

}
