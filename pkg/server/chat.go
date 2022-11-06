package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

// Client represents a machine sending and receiving messages
type Client chan<- string

var (
	incomingMsg = make(chan Client)
	outgoingMsg = make(chan Client)
	clientMsg   = make(chan string)

	host = flag.String("h", "localhost", "host")
	port = flag.Int("p", 1989, "port")
)

// handleConnection ...
func handleConnection(conn net.Conn) {
	defer conn.Close()
	// this is used to send general messages about users
	// joining and leaving the session
	systemMsg := make(chan string)
	go writeMsg(conn, systemMsg)
	// IP + Port
	clientName := conn.RemoteAddr().String()
	systemMsg <- fmt.Sprintf("Welcome to the server,%s\n", clientName)
	clientMsg <- fmt.Sprintf("New client is here, machine name: %s\n", clientName)
	incomingMsg <- systemMsg
	// scan the msg from the client
	inputMsg := bufio.NewScanner(conn)
	// it will loop indefinitely until the client disconnects (the connection is closed)
	for inputMsg.Scan() {
		clientMsg <- fmt.Sprintf("%s says: %s\n", clientName, inputMsg.Text())
	}
	// let others know about a client leaving the session
	outgoingMsg <- systemMsg
	clientMsg <- fmt.Sprintf("%s said goodbye!", clientName)
}

//goland:noinspection SpellCheckingInspection
func writeMsg(conn net.Conn, msgs <-chan string) {
	for msg := range msgs {
		_, _ = fmt.Fprintln(conn, msg)
	}
}

// broadcast Indefinitely listens to global channels and broadcast any new message.
func broadcast() {
	clients := make(map[Client]bool)
	for {
		select {
		case msg := <-clientMsg:
			for c := range clients {
				c <- msg
			}
		case newClient := <-incomingMsg:
			clients[newClient] = true
		case leavingClient := <-outgoingMsg:
			delete(clients, leavingClient)
			close(leavingClient)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatalln(err.Error())
	}
	go broadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

/*
The flow of the program is:
-> starts to listen on the given address using the tcp protocol.
-> fires a go routine that will be distributing incoming messages.
-> waits on the listener for any new connection and fires a goroutine to handle them accordingly.
*/
