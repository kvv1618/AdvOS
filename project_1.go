package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
	"time"
)

type JD struct {
	file_path string
	start_seg int
	seg_len   int
}

type partial_ans struct {
	JD
	num_primes int
}

func dispatcher(file_path string, n int, jobs_q chan JD, threads_g *sync.WaitGroup) {
	defer threads_g.Done()
	defer close(jobs_q)

	file, err := os.Open(file_path)
	defer file.Close()
	if err != nil {
		os.Exit(1)
	}

	segment := 0
	read_buf := make([]byte, n)
	var jd JD
	for {
		num_read_byes, err := file.Read(read_buf)
		if num_read_byes == 0 {
			break
		}
		if err == nil || err == io.EOF {
			jd = JD{file_path, segment, num_read_byes}
		} else {
			fmt.Println("Error reading file:\n", err)
			os.Exit(1)
		}
		jobs_q <- jd
		if err == io.EOF {
			break
		}
		segment += num_read_byes
	}
}

func is_prime(num int) bool {
	if num == 0 || num == 1 {
		return false
	}
	for i := 2; i <= num/2; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func nu_of_primes(read_buf []byte) int {
	byte_reader := bytes.NewReader(read_buf)
	var num uint64
	num_primes := 0
	for i := 0; i < len(read_buf); i++ {
		err := binary.Read(byte_reader, binary.LittleEndian, &num)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading from buffer:\n", err)
			os.Exit(1)
		}
		if err == io.EOF {
			break
		}
		if is_prime(int(num)) {
			num_primes++
		}
	}
	return num_primes
}

func worker(jobs_q chan JD, partial_ans_q chan partial_ans, c int, file_path string, threads_g *sync.WaitGroup, wg *sync.WaitGroup) {
	defer threads_g.Done()
	defer wg.Done()

	time.Sleep(time.Duration(rand.IntN(201)+400) * time.Millisecond)

	file, err := os.Open(file_path)
	defer file.Close()
	if err != nil {
		os.Exit(1)
	}

	for job := range jobs_q {
		file.Seek(int64(job.start_seg), 0)
		num_seg_to_read := 0
		if job.seg_len < c {
			num_seg_to_read = 1
		} else {
			num_seg_to_read = job.seg_len / c
		}
		for i := 0; i < num_seg_to_read; i++ {
			read_buf := make([]byte, c)
			num_read_bytes, err := file.Read(read_buf)
			if num_read_bytes == 0 {
				break
			}
			num_primes := nu_of_primes(read_buf)
			slog.Info("File from " + strconv.Itoa(job.start_seg+i*c) + " bytes to " + strconv.Itoa(job.start_seg+i*c+num_read_bytes) + " bytes has " + strconv.Itoa(num_primes) + " primes")
			partial_ans_q <- partial_ans{job, num_primes}
			if err == io.EOF {
				break
			}
		}
		if job.seg_len%c != 0 {
			read_buf := make([]byte, job.seg_len%c)
			num_read_byes, err := file.Read(read_buf)
			if num_read_byes == 0 {
				break
			}
			num_primes := nu_of_primes(read_buf)
			slog.Info("File from " + strconv.Itoa(job.start_seg+num_seg_to_read*c) + " bytes to " + strconv.Itoa(job.start_seg+job.seg_len) + " bytes has " + strconv.Itoa(num_primes) + " primes")
			partial_ans_q <- partial_ans{job, num_primes}
			if err == io.EOF {
				break
			}
		}
	}
}

func consolidator(partial_ans_q chan partial_ans, threads_g *sync.WaitGroup, num_primes *int) {
	defer threads_g.Done()
	for partial_ans := range partial_ans_q {
		*num_primes += partial_ans.num_primes
	}
}

func main() {
	file_path := os.Args[1]
	m, n, c := 1, 64000, 1000
	if len(os.Args) > 3 {
		_, err := fmt.Sscanf(os.Args[2], "%d", &m)
		if err != nil {
			fmt.Println("Defaulting value of m to 1")
		}
	}
	if len(os.Args) > 4 {
		_, err := fmt.Sscanf(os.Args[3], "%d", &n)
		if err != nil {
			fmt.Println("Defaulting value of n to 64000")
		}
	}
	if len(os.Args) > 5 {
		_, err := fmt.Sscanf(os.Args[4], "%d", &c)
		if err != nil {
			fmt.Println("Defaulting value of c to 1000")
		}
	}
	jobs_q, partial_ans_q := make(chan JD), make(chan partial_ans)
	num_primes := 0

	var wg, threads_g sync.WaitGroup

	threads_g.Add(1)
	go consolidator(partial_ans_q, &threads_g, &num_primes)

	for i := 0; i < m; i++ {
		wg.Add(1)
		threads_g.Add(1)
		go worker(jobs_q, partial_ans_q, c, file_path, &threads_g, &wg)
	}

	go func() {
		wg.Wait()
		defer close(partial_ans_q)
	}()

	threads_g.Add(1)
	go dispatcher(file_path, n, jobs_q, &threads_g)

	threads_g.Wait()

	fmt.Println(num_primes)
}
