package main

import (
	"fmt"
	"os"
)

func main() {
	file_path := os.Args[1]
	var m, n, c int
	_, err := fmt.Sscanf(os.Args[2], "%d", &m)
	_, err = fmt.Sscanf(os.Args[3], "%d", &n)
	_, err = fmt.Sscanf(os.Args[4], "%d", &c)
	file, err := os.Open(file_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	var num uint64

}
