package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/jakobilobi/go-wsstat"
)

func main() {
	// 1. Set up variables
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Usage: go run main.go URL")
	}
    wsURL := args[1]
	var err error

	// 2. Create a new WSStat instance
    ws := wsstat.NewWSStat()

	// 3. Establish a WebSocket connection
	// This triggers the DNS lookup, TCP connection, TLS handshake, and WebSocket handshake timers
	if err := ws.Dial(wsURL); err != nil {
		log.Fatalf("Failed to establish WebSocket connection: %v", err)
	}

	// 4. Send a message and wait for the response
	// This triggers the first message round trip timer
    // 4.a. Send a custom message and wait for the response
    msg := "Hello, WebSocket!"
	p, err := ws.SendMessage(websocket.TextMessage, []byte(msg))
    if err != nil {
        log.Fatalf("Failed to send message: %v", err)
    }
	log.Printf("Received message: %s", p)
	// 4.b. If not interested in the response, you can use the basic message sending method
	/* if err := ws.SendMessageBasic(); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	} */
	// 4.c. Alternatively just send a ping message
	/* if err := ws.SendPing(); err != nil {
		log.Fatalf("Failed to send ping: %v", err)
	} */

	// 5. Close the WebSocket connection
	// This triggers the total connection time timer, which includes the call to close the connection
	log.Println("Closing connection")
	err = ws.CloseConn()
	if err != nil {
		log.Fatalf("Failed to close WebSocket connection: %v", err)
	}

    // 6. Print the results
    fmt.Println("DNS Lookup Time:", ws.Result.DNSLookupDone)
    fmt.Println("TCP Connected:", ws.Result.TCPConnected)
    fmt.Println("TLS Handshake Done:", ws.Result.TLSHandshakeDone)
    fmt.Println("WS Handshake Done:", ws.Result.WSHandshakeDone)
    fmt.Println("First Message Round Trip Done:", ws.Result.FirstMessageResponse)
    fmt.Println("Total Connection Time:", ws.Result.TotalTime)
}
