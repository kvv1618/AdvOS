package main

import (
	"fmt"
	pb "github.com/kvv1618/Project2/protoc/service"
	"google.golang.org/grpc"
	"net"
	"os"
	"sync"
)

type fileServer struct {
	pb.UnimplementedJobDataServiceServer
}

func (s *fileServer) JobData(reg *pb.JobDetailsResponse, stream pb.JobDataService_JobDataServer) error {
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

	if err := stream.Send(&pb.JobDataResponse{Data: data}); err != nil {
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
	fmt.Printf("File server listening on port %s\n", ":5003")
	s := grpc.NewServer()
	pb.RegisterJobDataServiceServer(s, &fileServer{})
	if err := s.Serve(listner); err != nil {
		fmt.Println("Error serving gRPC server:\n", err)
		os.Exit(1)
	}
	defer s.Stop()
	defer listner.Close()
}

func main() {
	// data_file, config_file := os.Args[1], os.Args[2]
	data_file := os.Args[1]
	var wg sync.WaitGroup

	wg.Add(1)
	go fileserver(data_file, &wg)
	wg.Wait()
}
