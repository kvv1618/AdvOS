package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	pb "protoc/service"
	"strconv"
	"sync"
)

type JD struct {
	filePath string
	startSeg int
	segLen   int
}
type server struct {
	jobsQ       chan JD
	totalPrimes int
	pb.UnimplementedjobServiceServer
	pb.UnimplementedcondenseResultsServiceServer
	pb.UnimplementedjobDataServiceServer
}

func (s *server) jobDetails(ctx context.Context, req *pb.empty) (*pb.jobDetailsResponse, error) {
	jd := <-s.jobsQ
	if jd == (JD{}) {
		return nil, fmt.Errorf("no more jobs available")
	}
	resp := &pb.jobDetailsResponse{
		filePath: jd.filePath,
		startSeg: int32(jd.startSeg),
		segLen:   int32(jd.segLen),
	}
	return resp, nil
}

func dispatcher(filePath string, n int, c int, jobsQ chan JD, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(jobsQ)
	listner, err := net.Listen("tcp", ":5001")
	if err != nil {
		fmt.Println("Error listening on port 5001:\n", err)
		os.Exit(1)
	}
	defer listner.Close()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:\n", err)
		os.Exit(1)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:\n", err)
		os.Exit(1)
	}
	fileSize := fileInfo.Size()
	segment := 0
	for segment < int(fileSize) {
		jd := JD{filePath, segment, n}
		jobsQ <- jd
		segment += n
	}

	s := grpc.NewServer()
	pb.RegisterjobServiceServer(s, &server{
		jobsQ: jobsQ,
	})
	if err := s.Serve(listner); err != nil {
		fmt.Println("Error serving gRPC server:\n", err)
		os.Exit(1)
	}

	select {}
}

func (s *server) condenseResults(ctx context.Context, req *pb.partialResults) (*pb.empty, error) {
	s.totalPrimes += int(req.NumPrimes)
	return &pb.empty{}, nil
}

func consolidator(wg *sync.WaitGroup) {
	defer wg.Done()
	listner, err := net.Listen("tcp", ":5002")
	if err != nil {
		fmt.Println("Error listening on port 5002:\n", err)
		os.Exit(1)
	}
	defer listner.Close()

	s := grpc.NewServer()
	pb.RegisterPartialResultsServiceServer(s, &server{})
	if err := s.Serve(listner); err != nil {
		fmt.Println("Error serving gRPC server:\n", err)
		os.Exit(1)
	}

	select {}
}

func main() {
	N, _ := strconv.Atoi(os.Args[1])
	C, _ := strconv.Atoi(os.Args[2])
	data_file, config_file := os.Args[3], os.Args[4]

	jobsQ := make(chan JD)

	var wg sync.WaitGroup
	wg.Add(1)
	go dispatcher(data_file, N, C, jobsQ, &wg)

	wg.Add(1)
	go consolidator(&wg)

	wg.Wait()
}
