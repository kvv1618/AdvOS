package main

import (
	"fmt"
	"os"
)

// worker reads each segment C bytes at a time
func main() {
	C, _ := strconv.Atoi(os.Args[1])
	config_file := os.Args[2]

}
