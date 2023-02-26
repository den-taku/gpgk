package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
	"os/exec"
	"fmt"
	"github.com/google/uuid"
	pb "server/grpc"
)

type server struct {
	pb.UnimplementedGpgkServiceServer
}

func writeFile(filename, code string)  {
	f, err := os.Create(filename + ".rs")
	if err != nil {
		log.Fatalf("failed to open %v: %v", filename + ".rs", err)
		return
	}
	defer f.Close()
	count, err := f.Write([]byte(code))
	if err != nil {
		log.Fatalf("failed to write: %v", err)
	}
	log.Printf("new file created, %v bytes", count)
}

func executeFile(filename string) string {
	out, err := exec.Command("/Users/dentaku/.cargo/bin/rustc", filename + ".rs").Output()
	if err != nil {
		log.Fatalf("failed to compile: %v, %v", out, err)
	}
	log.Printf("compile: %v", out)
	stdout, err := exec.Command("./" + filename).Output()
	if err != nil {
		log.Fatalf("error occurs: %v", err)
	}
	return string(stdout)
}

func removeFiles(filename string) {
	// move .rs to log
	out1, err := exec.Command("mv", filename + ".rs", "log/" + filename + ".rs").Output()
	if err != nil {
		log.Fatalf("failed to move file: %v", err)
	}
	log.Printf("move file to log: %v", out1)

	// remove elf
	out2, err := exec.Command("rm", filename).Output()
	if err != nil {
		log.Fatalf("failed to remove file: %v", err)
	}
	log.Printf("remove file: %v", out2)
}

func (s *server) Execute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	code := req.GetCode()
	log.Printf("RECV: %v", code)

	// make unique filename with time
	t := time.Now().Unix()
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("failed to generate uuid: %v", err)
	}
	filename := "code" + fmt.Sprintf("%d", t) + "_" + uuid.String()

	// exec
	writeFile(filename, code)
	stdout := executeFile(filename)
	removeFiles(filename)

	log.Printf("SEND: %v", stdout)
	return &pb.ExecuteResponse{Stdout: stdout}, nil
}

func main() {
	// create log
	out, err := exec.Command("mkdir", "-p", "log").Output()
	if err != nil {
		log.Fatalf("Failed to create log: %v", err)
		return
	}
	log.Printf("create log: %v", out)

	// wake up server
	lis, err := net.Listen("tcp", ":5859")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGpgkServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
