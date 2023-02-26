package main

import (
	"context"
	"log"
	"time"
	"google.golang.org/grpc"
	pb "server/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:5959", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewGpgkServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()
	r, err := c.Execute(ctx, &pb.ExecuteRequest{Code: "fn main() {println!(\"Hello!!\")}"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("RECV: %s", r.GetStdout())
}