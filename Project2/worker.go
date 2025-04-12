package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"math/rand/v2"
	"os"
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

func worker(C int, config_file string, wg *sync.WaitGroup) {
	defer wg.Done()
	dispatcherConn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to dispatcher:\n", err)
		os.Exit(1)
	}
	defer dispatcherConn.Close()

	consolidatorConn, err := grpc.NewClient("localhost:5002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to consolidator:\n", err)
		os.Exit(1)
	}
	defer consolidatorConn.Close()

	fileServerConn, err := grpc.NewClient("localhost:5003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to file server:\n", err)
		os.Exit(1)
	}
	defer fileServerConn.Close()

	dispatcherClient := pb.NewJobServiceClient(dispatcherConn)
	fileServerClinet := pb.NewJobDataServiceClient(fileServerConn)
	consolidaterClient := pb.NewCondenseResultsServiceClient(consolidatorConn)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		resp, err := dispatcherClient.JobDetails(ctx, &pb.Empty{})
		if err != nil && strings.Contains(err.Error(), "no more jobs available") {
			break
		}
		if err != nil || resp == nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cancel()

		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
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
		}
		cancel()

		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
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

		time.Sleep(time.Duration(rand.IntN(201)+400) * time.Millisecond)
	}
}

func main() {
	m, C := 1, 1000
	if len(os.Args) > 1 {
		C, _ = strconv.Atoi(os.Args[1])
	}
	// config_file := os.Args[2]
	if len(os.Args) > 3 {
		m, _ = strconv.Atoi(os.Args[3])
	}

	var wg sync.WaitGroup
	wg.Add(1)
	for i := 0; i < m; i++ {
		// go worker(C, config_file)
		go worker(C, "", &wg)
	}
	wg.Wait()

	conn, err := grpc.NewClient("localhost:5002", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
