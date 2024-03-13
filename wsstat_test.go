package wsstat

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var (	
	serverAddr = "localhost:8080"
	echoServerAddrWs = "ws://" + serverAddr + "/echo"
)

// TestMain sets up the test server and runs the tests in this file.
func TestMain(m *testing.M) {
	// Set up test server
	go func() {
		StartEchoServer(serverAddr)
	}()

	// Ensure the echo server starts before any tests run
	time.Sleep(250 * time.Millisecond)

	// Run the tests in this file
	os.Exit(m.Run())
}

func TestMeasureLatency(t *testing.T) {
	msg := "Hello, world!"
	result, response, err := MeasureLatency(echoServerAddrWs, msg)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.TotalTime <= 0 {
		t.Errorf("Invalid total time: %v", result.TotalTime)
	}
	if len(response) == 0 {
		t.Errorf("Empty response")
	}
	if string(response) != msg {
		t.Errorf("Unexpected response: %s", response)
	}
}

func TestMeasureLatencyJSON(t *testing.T) {
	message := struct {
		Text string `json:"text"`
	}{
		Text: "Hello, world!",
	}
	result, response, err := MeasureLatencyJSON(echoServerAddrWs, message)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.TotalTime <= 0 {
		t.Errorf("Invalid total time: %v", result.TotalTime)
	}
	if response == nil {
		t.Errorf("Empty response")
	}
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		t.Errorf("Response is not a map")
		return
	}
	if responseMap["text"] != message.Text {
		t.Errorf("Unexpected response: %s", responseMap)
	}
}

// Helpers

// StartEchoServer starts a WebSocket server that echoes back any received messages.
func StartEchoServer(addr string) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins
		},
	}

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer conn.Close()
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				// Only print error if it's not a normal closure
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("read:", err)
				}
				break
			}
			err = conn.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	})

	fmt.Printf("Echo server started on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
