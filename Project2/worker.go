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
	pb "protoc/service"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
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

func worker(C int, config_file string) {
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Error connecting to dispatcher:\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	dispatcherClient := pb.NewjobServiceClient(conn)
	fileServerClinet := pb.NewjobDataServiceClient(conn)
	consolidaterClient := pb.NewcondenseResultsServiceClient(conn)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := dispatcherClient.jobDetails(ctx, &pb.empty{})
		cancel()
		if err != nil {
			fmt.Println("Error getting job details:\n", err)
			break
		}
		if resp == nil {
			break
		}

		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		stream, err := fileServerClinet.jobData(ctx, &pb.jobDetailsResponse{
			FilePath: resp.FilePath,
			StartSeg: resp.StartSeg,
			segLen:   resp.segLen,
		})
		if err != nil {
			fmt.Println("Error getting job data:\n", err)
			break
		}
		readBuf := make([]byte, C)
		numPrimes := 0
		for {
			numReadBytes, err := stream.RecvMsg(&pb.jobDataResponse{data: readBuf})
			if numReadBytes == 0 {
				break
			}
			numPrimes += findPrimes(readBuf)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Error reading stream data:\n", err)
				break
			}
		}

		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		partialAns := &pb.partialResults{
			FilePath:  resp.FilePath,
			StartSeg:  resp.StartSeg,
			segLen:    resp.segLen,
			numPrimes: int32(numPrimes),
		}
		_, err = consolidaterClient.condenseResults(ctx, partialAns)
		cancel()
		if err != nil {
			fmt.Println("Error sending partial results:\n", err)
			break
		}

		time.Sleep(time.Duration(rand.IntN(201)+400) * time.Millisecond)
	}
}

func main() {
	m := 1
	C, _ := strconv.Atoi(os.Args[1])
	config_file := os.Args[2]
	if len(os.Args) > 3 {
		m, _ = strconv.Atoi(os.Args[3])
	}

	var wg sync.WaitGroup
	wg.Add(1)

	for i := 0; i < m; i++ {
		go worker(C, config_file)
	}

}
