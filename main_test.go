package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestGracefulShutdownSuccess(t *testing.T) {
	// Set up a buffer to capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr) // Restore default output

	router := CreateRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	// Use the test server's URL to ensure we are connecting to the correct port
	port := server.Listener.Addr().(*net.TCPAddr).Port
	t.Logf("Server started on: %s", server.URL)

	go func() {
		// Create the JSON body with the sleep value
		jsonBody := `{"sleep": 5}`
		req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/item/123", port), strings.NewReader(jsonBody))
		if err != nil {
			log.Printf("Error creating POST request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// Perform the POST request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error making POST request: %v", err)
			return
		}
		defer resp.Body.Close()

		log.Printf("POST request successful: %v", resp.Status)
	}()

	// Simulate an OS interrupt after a short delay
	go func() {
		time.Sleep(3 * time.Second)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGINT)
	}()

	GracefulShutdown(server.Config, 10*time.Second)

	// Check the log output for the expected message
	logOutput := logBuffer.String()
	notExpectedMessage := "Could not gracefully shut down the server:"
	if bytes.Contains([]byte(logOutput), []byte(notExpectedMessage)) {
		t.Errorf("Unexpected log message found: %s; last received message: %s", notExpectedMessage, logOutput)
	}
}

func TestGracefulShutdownTimeout(t *testing.T) {
	// Set up a buffer to capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr) // Restore default output

	router := CreateRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	// Get the port number from the returned server URL
	port := server.Listener.Addr().(*net.TCPAddr).Port
	t.Logf("Server started on: %s", server.URL)

	// Simulate an HTTP POST request to /item/123
	go func() {
		// Create the JSON body with the sleep value
		jsonBody := `{"sleep": 5}`
		req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/item/123", port), strings.NewReader(jsonBody))
		if err != nil {
			log.Printf("Error creating POST request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// Perform the POST request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error making POST request: %v", err)
			return
		}
		defer resp.Body.Close()

		log.Printf("POST request successful: %v", resp.Status)
	}()

	// Simulate an OS interrupt after a short delay
	go func() {
		time.Sleep(5 * time.Second)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGINT)
	}()

	// Perform graceful shutdown with a very short timeout
	GracefulShutdown(server.Config, 1*time.Millisecond)

	// Check the log output for the expected message
	logOutput := logBuffer.String()
	expectedMessage := "Could not gracefully shut down the server:"
	if !bytes.Contains([]byte(logOutput), []byte(expectedMessage)) {
		t.Errorf("Expected log message not found: %s; last received message: %s", expectedMessage, logOutput)
	}
}
