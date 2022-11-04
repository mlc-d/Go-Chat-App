package main

import (
	"flag"
	"fmt"
)

var (
	host = flag.String("host", "localhost", "")
	port = flag.String("port", "1989", "")
)

// port scan
func main() {
	fmt.Println("hey!")
}
