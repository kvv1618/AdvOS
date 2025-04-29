package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/kvv1618/Project2/protoc/service"
	"google.golang.org/grpc"
)

type JD struct {
	filePath string
	startSeg int
	segLen   int
}
type server struct {
	startTime   time.Time
	jobsQ       chan JD
	totalPrimes int
	pb.UnimplementedJobServiceServer
	pb.UnimplementedCondenseResultsServiceServer
	pb.UnimplementedStopConsolidatorServiceServer
}

func (s *server) JobDetails(ctx context.Context, req *pb.Empty) (*pb.JobDetailsResponse, error) {
	select {
	case jd := <-s.jobsQ:
		resp := &pb.JobDetailsResponse{
			FilePath: jd.filePath,
			StartSeg: int64(jd.startSeg),
			SegLen:   int64(jd.segLen),
		}
		return resp, nil
	default:
		return nil, fmt.Errorf("no more jobs available")
	}
}

func dispatcher(filePath string, n int, c int, wg *sync.WaitGroup, str_port string) {
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
	defer listner.Close()
	fmt.Printf("Dispatcher listening on port %s\n", fmt.Sprintf(":%d", port))

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fileStats, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file stats:\n", err)
		os.Exit(1)
	}
	channelMaxSize := int(math.Ceil(float64(fileStats.Size()) / float64(n)))
	jobsQ := make(chan JD, channelMaxSize)

	segment := 0
	read_buf := make([]byte, n)
	for {
		numReadBytes, err := file.Read(read_buf)
		if numReadBytes == 0 {
			break
		}
		if err == nil || err == io.EOF {
			jd := JD{filePath, segment, numReadBytes}
			jobsQ <- jd
		} else {
			fmt.Println("Error reading file:\n", err)
			os.Exit(1)
		}
		if err == io.EOF {
			break
		}
		segment += numReadBytes
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
	go func() {
		fmt.Println("Elapsed time (ms):", time.Since(s.startTime).Milliseconds(), "ms")

		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()
	return &pb.Empty{}, nil
}

func consolidator(wg *sync.WaitGroup, str_port string, start time.Time) {
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
	fmt.Printf("Consolidator listening on port %s\n", fmt.Sprintf(":%d", port))
	s := grpc.NewServer()
	serverInstance := &server{
		startTime: start,
	}
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
	start := time.Now()

	N, C := 64000, 1000
	if len(os.Args) > 1 {
		N, _ = strconv.Atoi(os.Args[1])
	}
	if len(os.Args) > 2 {
		C, _ = strconv.Atoi(os.Args[2])
	}
	data_file, config_file := os.Args[3], os.Args[4]
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
	go dispatcher(data_file, N, C, &wg, ports["dispatcher"]["port"])

	wg.Add(1)
	go consolidator(&wg, ports["consolidator"]["port"], start)

	wg.Wait()

}
