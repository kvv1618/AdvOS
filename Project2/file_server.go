package main

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	pb "protoc/service"
	"sync"
)

func (s *server) jobData(reg *pb.jobDetailsResponse, stream pb.jobDataService_jobDataServer) error {
	startSeg := reg.StartSeg
	segLen := reg.SegLen
	filePath := reg.FilePath

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	if _, err := file.Seek(int64(startSeg), 0); err != nil {
		return fmt.Errorf("failed to seek to start segment: %v", err)
	}
	data := make([]byte, segLen)
	if _, err := file.Read(data); err != nil {
		return fmt.Errorf("failed to read segment: %v", err)
	}

	if err := stream.Send(&pb.jobDataResponse{data: data}); err != nil {
		return fmt.Errorf("failed to send segment data: %v", err)
	}

	return nil
}

func fileserver(data_file string, wg *sync.WaitGroup) {
	defer wg.Done()
	listner, err := net.Listen("tcp", ":5003")
	if err != nil {
		fmt.Println("Error listening on port 5003:\n", err)
		os.Exit(1)
	}
	defer listner.Close()

	s := grpc.NewServer()
	pb.RegisterjobDataServiceServer(s, &server{})

	if err := s.Serve(listner); err != nil {
		fmt.Println("Error serving gRPC server:\n", err)
		os.Exit(1)
	}

	select {}
}

func main() {
	data_file, config_file := os.Args[1], os.Args[2]
	var wg sync.WaitGroup

	wg.Add(1)
	go fileserver(data_file, &wg)
	wg.Wait()
}
