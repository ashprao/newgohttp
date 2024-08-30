package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body><h1>404 - Not Found</h1><p>The requested URL %s was not found on this server.</p></body></html>", r.URL.Path)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><h1>Hello, World!</h1></body></html>")
}

func itemIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><h1>Item: %s</h1></body></html>", id)
}

// Struct to represent the expected JSON body
type PostRequestBody struct {
	Sleep int `json:"sleep,omitempty"`
}

func itemPostHandler(w http.ResponseWriter, r *http.Request) {
	// Read the body of the POST request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Initialize default values
	sleepDuration := 0

	// Parse the JSON body if it exists
	if len(body) > 0 {
		var reqBody PostRequestBody
		if err := json.Unmarshal(body, &reqBody); err != nil {
			http.Error(w, "Invalid request body format", http.StatusBadRequest)
			return
		}

		// Set the sleep duration from the request
		sleepDuration = reqBody.Sleep
	}

	// Sleep for the specified duration if greater than 0
	if sleepDuration > 0 {
		log.Printf("Sleeping for %d seconds", sleepDuration)
		time.Sleep(time.Duration(sleepDuration) * time.Second)
	}

	// Extract the item ID from the URL
	id := r.URL.Path[len("/item/"):]

	// Send an HTML response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><h1>Post Item: %s</h1></body></html>", id)
}

// CreateRouter initializes the HTTP routes and returns a ServeMux
func CreateRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("GET /item/{id}", itemIdHandler)
	router.HandleFunc("POST /item/{id}", itemPostHandler)
	return router
}

// StartServer starts the HTTP server with the provided router and address
func StartServer(addr string, router http.Handler) *http.Server {
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Starting server in a goroutine to allow graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()
	log.Printf("Server is ready to handle requests at http://localhost%s", server.Addr)

	return server
}

// GracefulShutdown handles the shutdown process of the server
func GracefulShutdown(server *http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	log.Printf("Server at %s shutting down...", server.Addr)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		// Check if the context's deadline was exceeded, indicating a timeout
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Graceful shutdown timed out, forcing exit.")
		}
		log.Printf("Could not gracefully shut down the server: %v\n", err)
	} else {
		log.Printf("Server at %s stopped\n", server.Addr)
	}

}

func main() {
	router := CreateRouter()
	server := StartServer(":8080", router)
	GracefulShutdown(server, 10*time.Second)
}
