package main

import (
	"context"
	"log"
	"time"
	"google.golang.org/grpc"
	pb "server/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:5859", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewGpgkServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()
	program := `// this code is sent from test-client
fn main() {
	println!("Hello, DenTaku!")
}
`
	r, err := c.Execute(ctx, &pb.ExecuteRequest{Code: program})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("RECV: %s", r.GetStdout())
}