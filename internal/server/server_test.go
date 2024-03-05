package server

import (
	"net"
	"sync"
	"testing"
	"time"
)

const (
	serverAddr = ":6120"
	numClients = 15
)

// creates and starts a test server and returns its instance.
func startTestServer(t *testing.T, serverAddr string) *Server {
	server, err := NewServer(serverAddr)
	if err != nil {
		t.Fatalf("Error creating server: %v", err)
	}
	go func() {
		err := server.StartServer()
		if err != nil {
			t.Errorf("Server error: %v", err)
		}
	}()
	// Give the server some time to start
	time.Sleep(100 * time.Millisecond)
	return server
}

// creates a test client that connects to the server, sends a query, and validates the response.
func createTestClient(t *testing.T, serverAddr string, clientNum int, responses chan bool) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Errorf("Error connecting to server (client %d): %v", clientNum, err)
		return
	}

	query := "pwd\n"
	_, err = conn.Write([]byte(query))
	if err != nil {
		t.Errorf("Error writing to server (client %d): %v", clientNum, err)
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Errorf("Error reading from server (client %d): %v", clientNum, err)
		return
	}

	response := string(buffer[:n])
	expectedResponse := "Error: query must start with 'grep'"
	if response != expectedResponse {
		t.Errorf("Expected response '%s', but got '%s' (client %d)", expectedResponse, response, clientNum)
		responses <- false
	} else {
		responses <- true
	}
}

// TestServerMultipleClients tests the server with multiple concurrent clients.
func TestServerMultipleClients(t *testing.T) {
	startTestServer(t, serverAddr)

	var wg sync.WaitGroup
	wg.Add(numClients)

	responses := make(chan bool, numClients)

	for i := 0; i < numClients; i++ {
		go func(clientNum int) {
			defer wg.Done()
			createTestClient(t, serverAddr, clientNum, responses)
		}(i)
	}
	wg.Wait()

	if len(responses) != numClients {
		t.Errorf("Expected %d responses, but received %d", numClients, len(responses))
	}

}
