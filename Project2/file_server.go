package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	pb "github.com/kvv1618/Project2/protoc/service"
	"google.golang.org/grpc"
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

func fileserver(data_file string, wg *sync.WaitGroup, str_port string) {
	defer wg.Done()
	port, err := strconv.Atoi(str_port)
	if err != nil {
		fmt.Println("Error converting port to int:\n", err)
		os.Exit(1)
	}

	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error listening on port %d:\n", port, err)
		os.Exit(1)
	}
	fmt.Printf("File server listening on port %s\n", fmt.Sprintf(":%d", port))
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
	data_file, config_file := os.Args[1], os.Args[2]
	file, err := os.Open(config_file)
	if err != nil {
		fmt.Println("Error opening config file:\n", err)
		os.Exit(1)
	}
	defer file.Close()

	ports := make(map[string]map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		if len(line) != 3 {
			fmt.Println("Invalid config file format")
			os.Exit(1)
		}
		ports[line[0]] = map[string]string{
			"port": line[2],
			"ip":   line[1],
		}
	}
	var wg sync.WaitGroup

	wg.Add(1)
	go fileserver(data_file, &wg, ports["fileserver"]["port"])
	wg.Wait()
}
