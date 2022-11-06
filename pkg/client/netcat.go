package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	port = flag.Int("port", 1989, "port")
	host = flag.String("host", "localhost", "host")
)

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatalln(err.Error())
	}
	// simple channel used to block the execution of the program
	done := make(chan struct{})
	go func() {
		copyContent(os.Stdout, conn)
		done <- struct{}{}
	}()
	copyContent(conn, os.Stdin)
	_ = conn.Close()
	<-done
}

// copyContent envelope for [io.Copy], which panics if the operation fails.
// Used to copy to and from the system standard output and input, respectively.
func copyContent(w io.Writer, r io.Reader) {
	_, err := io.Copy(w, r)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

/*
The flow of the program is:
-> It dials the given tcp address (line 26).
-> it fires a go routine that is copying the content from the tcp connection to stdout.
-> the main routine then starts to read from the stdin and writing to the tcp connection.
NOTE: as far as i can see, lines 30 and 31 are never executed.
*/
