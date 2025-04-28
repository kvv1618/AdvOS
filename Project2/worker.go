package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"math/rand/v2"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/kvv1618/Project2/protoc/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func isPrime(num uint64) bool {
	bigNum := new(big.Int).SetUint64(uint64(num))
	return bigNum.ProbablyPrime(10)
}

func findPrimes(readBuf []byte) int {
	byte_reader := bytes.NewReader(readBuf)
	var num uint64
	num_primes := 0
	for i := 0; i < len(readBuf); i++ {
		err := binary.Read(byte_reader, binary.LittleEndian, &num)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading from buffer:\n", err)
			os.Exit(1)
		}
		if err == io.EOF {
			break
		}
		if isPrime(uint64(num)) {
			num_primes++
		}
	}
	return num_primes
}

func worker(C int, wg *sync.WaitGroup, ports map[string]map[string]string, workerStats *sync.Map, id int) {
	defer wg.Done()

	dispatcherPort, err := strconv.Atoi(ports["dispatcher"]["port"])
	if err != nil {
		fmt.Println("Error converting dispatcher port to int:\n", err)
		os.Exit(1)
	}
	dispatcherIp := ports["dispatcher"]["ip"]
	dispatcherConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", dispatcherIp, dispatcherPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to dispatcher:\n", err)
		os.Exit(1)
	}
	defer dispatcherConn.Close()

	consolidatorPort, err := strconv.Atoi(ports["consolidator"]["port"])
	if err != nil {
		fmt.Println("Error converting consolidator port to int:\n", err)
		os.Exit(1)
	}
	consolidatorIp := ports["consolidator"]["ip"]
	consolidatorConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", consolidatorIp, consolidatorPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to consolidator:\n", err)
		os.Exit(1)
	}
	defer consolidatorConn.Close()

	fileServerPort, err := strconv.Atoi(ports["fileserver"]["port"])
	if err != nil {
		fmt.Println("Error converting file server port to int:\n", err)
		os.Exit(1)
	}
	fileServerIp := ports["fileserver"]["ip"]
	fileServerConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", fileServerIp, fileServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to file server:\n", err)
		os.Exit(1)
	}
	defer fileServerConn.Close()

	dispatcherClient := pb.NewJobServiceClient(dispatcherConn)
	fileServerClinet := pb.NewJobDataServiceClient(fileServerConn)
	consolidaterClient := pb.NewCondenseResultsServiceClient(consolidatorConn)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := dispatcherClient.JobDetails(ctx, &pb.Empty{})
		if err != nil && strings.Contains(err.Error(), "no more jobs available") {
			break
		}
		if err != nil || resp == nil {
			// fmt.Println("Error getting job details:\n", err)
			cancel()
			time.Sleep(time.Duration(rand.IntN(201)+400) * time.Millisecond)
			continue
		}
		cancel()

		ctx, cancel = context.WithCancel(context.Background())
		stream, err := fileServerClinet.JobData(ctx, &pb.JobDetailsResponse{
			FilePath: resp.FilePath,
			StartSeg: resp.StartSeg,
			SegLen:   resp.SegLen,
		})
		if err != nil {
			fmt.Println("Error getting job data:\n", err)
			os.Exit(1)
		}

		numPrimes := 0
		numReadBytes, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading stream data:\n", err)
			os.Exit(1)
		}
		for i := 0; i < len(numReadBytes.Data); i += C {
			end := i + C
			if end > len(numReadBytes.Data) {
				end = len(numReadBytes.Data)
			}
			readBuf := numReadBytes.Data[i:end]
			numPrimes += findPrimes(readBuf)
			numPrimes += 0
		}
		cancel()

		ctx, cancel = context.WithCancel(context.Background())
		partialAns := &pb.PartialResults{
			FilePath:  resp.FilePath,
			StartSeg:  resp.StartSeg,
			SegLen:    resp.SegLen,
			NumPrimes: int32(numPrimes),
		}
		_, err = consolidaterClient.CondenseResults(ctx, partialAns)
		if err != nil {
			fmt.Println("Error sending partial results:\n", err)
			os.Exit(1)
		}
		cancel()

		if val, ok := workerStats.Load(id); ok {
			workerStats.Store(id, val.(int)+1)
		} else {
			workerStats.Store(id, 1)
		}
		time.Sleep(time.Duration(rand.IntN(201)+400) * time.Millisecond)
	}
}

func main() {
	m, C := 1, 1000
	if len(os.Args) > 1 {
		m, _ = strconv.Atoi(os.Args[1])
	}
	config_file := os.Args[2]
	if len(os.Args) > 3 {
		C, _ = strconv.Atoi(os.Args[3])
	}

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
	var workerStats sync.Map
	for i := 0; i < m; i++ {
		wg.Add(1)
		go worker(C, &wg, ports, &workerStats, i)
	}
	wg.Wait()

	jobStats := make([]int, 0, m)
	workerStats.Range(func(key, value interface{}) bool {
		jobStats = append(jobStats, value.(int))
		return true
	})
	sort.Ints(jobStats)
	min, max := jobStats[0], jobStats[len(jobStats)-1]
	sum, median := 0, 0
	for _, val := range jobStats {
		sum += val
	}
	average := float64(sum) / float64(len(jobStats))
	if len(jobStats)%2 == 0 {
		median = (jobStats[len(jobStats)/2-1] + jobStats[len(jobStats)/2]) / 2
	} else {
		median = jobStats[len(jobStats)/2]
	}
	fmt.Printf("Total jobs completed: %d\n", sum)
	fmt.Printf("Worker stats(Jobs completed by a worker):\n")
	fmt.Printf("Min: %d, Max: %d, Average: %.2f, Median: %d\n", min, max, average, median)

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", ports["consolidator"]["ip"], ports["consolidator"]["port"]), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to consolidator:\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	consolidaterClient := pb.NewStopConsolidatorServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_, err = consolidaterClient.StopConsolidator(ctx, &pb.Empty{})
	if err != nil {
		fmt.Println("Error stopping consolidator:\n", err)
		os.Exit(1)
	}
	cancel()
}
