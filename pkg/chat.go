package pkg

// Client represents the machine sending and receiving messages
type Client chan<- string

var (
	incomingMsg = make(chan Client)
	outgoingMsg = make(chan Client)
	messages    = make(chan string)
)
