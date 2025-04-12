package main

import (
	"container/list"
	"context"
	"fmt"
	pb "github.com/kvv1618/Project2/protoc/service"
	"google.golang.org/grpc"
	"net"
	"os"
	"strconv"
	"sync"
)

type JD struct {
	filePath string
	startSeg int
	segLen   int
}
type server struct {
	jobsQ       *list.List
	totalPrimes int
	pb.UnimplementedJobServiceServer
	pb.UnimplementedCondenseResultsServiceServer
	pb.UnimplementedStopConsolidatorServiceServer
}

func (s *server) JobDetails(ctx context.Context, req *pb.Empty) (*pb.JobDetailsResponse, error) {
	if s.jobsQ.Len() == 0 {
		return nil, fmt.Errorf("no more jobs available")
	}
	jdElement := s.jobsQ.Front()
	jd, ok := jdElement.Value.(JD)
	if !ok {
		return nil, fmt.Errorf("failed to cast job details")
	}
	resp := &pb.JobDetailsResponse{
		FilePath: jd.filePath,
		StartSeg: int32(jd.startSeg),
		SegLen:   int32(jd.segLen),
	}
	s.jobsQ.Remove(jdElement)
	return resp, nil
}

func dispatcher(filePath string, n int, c int, jobsQ *list.List, wg *sync.WaitGroup) {
	defer wg.Done()
	listner, err := net.Listen("tcp", ":5001")
	if err != nil {
		fmt.Println("Error listening on port 5001:\n", err)
		os.Exit(1)
	}
	defer listner.Close()
	fmt.Printf("Dispatcher listening on port %s\n", ":5001")

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
		jobsQ.PushBack(jd)
		segment += n
	}

	s := grpc.NewServer()
	pb.RegisterJobServiceServer(s, &server{
		jobsQ: jobsQ,
	})
	if err := s.Serve(listner); err != nil {
		fmt.Println("Error serving gRPC server:\n", err)
		os.Exit(1)
	}
	defer s.Stop()
}

func (s *server) CondenseResults(ctx context.Context, req *pb.PartialResults) (*pb.Empty, error) {
	s.totalPrimes += int(req.NumPrimes)
	fmt.Printf("Received partial results: %s, %d, %d, %d\n", req.FilePath, req.StartSeg, req.SegLen, req.NumPrimes)
	return &pb.Empty{}, nil
}

func (s *server) StopConsolidator(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	fmt.Println("Total number of primes found:", s.totalPrimes)
	os.Exit(0)
	return &pb.Empty{}, nil
}

func consolidator(wg *sync.WaitGroup) {
	defer wg.Done()
	listner, err := net.Listen("tcp", ":5002")
	if err != nil {
		fmt.Println("Error listening on port 5002:\n", err)
		os.Exit(1)
	}
	fmt.Printf("Consolidator listening on port %s\n", ":5002")
	s := grpc.NewServer()
	serverInstance := &server{}
	pb.RegisterCondenseResultsServiceServer(s, serverInstance)
	pb.RegisterStopConsolidatorServiceServer(s, serverInstance)
	if err := s.Serve(listner); err != nil {
		fmt.Println("Error serving gRPC server:\n", err)
		os.Exit(1)
	}
	defer s.Stop()
	defer listner.Close()
}

func main() {
	N, C := 64000, 1000
	if len(os.Args) > 1 {
		N, _ = strconv.Atoi(os.Args[1])
	}
	if len(os.Args) > 2 {
		C, _ = strconv.Atoi(os.Args[2])
	}
	// data_file, config_file := os.Args[3], os.Args[4]
	data_file := os.Args[3]

	jobsQ := list.New()

	var wg sync.WaitGroup
	wg.Add(1)
	go dispatcher(data_file, N, C, jobsQ, &wg)

	wg.Add(1)
	go consolidator(&wg)

	wg.Wait()

}
