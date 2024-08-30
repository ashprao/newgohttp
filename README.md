# Go HTTP Server Example

This Go HTTP server example demonstrates not only basic routing, handling GET and POST requests but also the implementation of a graceful shutdown process. 

It showcases the new features introduced in Go 1.22's `net/http` package, which have made path-based routing much more straightforward and powerful. Previously, such routing required cumbersome workarounds, prompting many developers to turn to external frameworks like Gin, Gorilla or Chi. With the advancements in Go 1.22, you can now achieve clean and efficient routing directly within the standard library, simplifying development and reducing dependencies. 

Additionally, this example demonstrates how to implement a graceful shutdown, allowing the server to handle in-flight requests properly before shutting down, ensuring a smooth and reliable service termination. This feature is crucial for maintaining stability and ensuring that active connections are not abruptly cut off during server shutdowns.

Support for adding middleware has also been incorporated in the `net/http` package. This example will be progressively updated to showcase these other capabilities.

## Features

- **Home Handler**: Responds with a "Hello, World!" message on the root path `/`. If the path value is not empty (i.e., if the user tries to access an unhandled path like `/unknown`), it returns a "404 - Not Found" page.
- **Item ID Handler**: Responds with the item ID when accessed via the path `/item/{id}`.
- **Item POST Handler**: Accepts a JSON body containing a `sleep` value to simulate a long-running task. It then responds with the item ID after the specified sleep duration. If no sleep value is provided, it responds immediately.
- **Graceful Shutdown**: The server supports a graceful shutdown mechanism to properly handle in-flight requests before shutting down.

## Usage

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.22 or later)

### Running the Server

1. Clone this repository:

   ```bash
   git clone https://github.com/yourusername/go-http-server-example.git
   cd go-http-server-example
   ```

2. Run the server:

   ```bash
   go run main.go
   ```

3. The server will start and listen on port `8080`. You can access it at `http://localhost:8080`.

### Endpoints

- **GET `/`**: Returns a simple HTML page with "Hello, World!". If any additional path is provided (e.g., `/unknown`), it returns a "404 - Not Found" page.

  Example:
  ```bash
  curl http://localhost:8080/
  ```

  For an unhandled path:

  ```bash
  curl http://localhost:8080/unknown
  ```

- **GET `/item/{id}`**: Returns a simple HTML page displaying the item ID.

  Example:
  ```bash
  curl http://localhost:8080/item/123
  ```

- **POST `/item/{id}`**: Accepts a JSON body containing a `sleep` value in seconds. The server sleeps for the specified duration before responding with the item ID. If no `sleep` value is provided, the server responds immediately.

  Example:
  ```bash
  curl -X POST http://localhost:8080/item/123 -H "Content-Type: application/json" -d '{"sleep": 15}'
  ```

  This example sends a POST request to `/item/123`, causing the server to sleep for 15 seconds before responding.

### Graceful Shutdown

- **Initiation**: The server supports graceful shutdown, allowing it to terminate connections and stop accepting new requests while finishing in-flight requests. To trigger a graceful shutdown, you can interrupt the process with `Ctrl+C`.

  A channel is created to receive operating system signals (like `ctrl+c` or terminal interrupts). When a signal is received (`<-c`), it initiates a shutdown process.

  ```go
    c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
  ```

- **Shutdown Context**: A context with a timeout is created using `context.WithTimeout`. This context is used to control the graceful shutdown process and ensure that it doesn't take too long.

- **Stopping the Server**: The `server.Shutdown(ctx)` function is called, which gracefully stops the server. It waits for any ongoing requests to complete (up to the specified timeout) before shutting down the server.

## Unit Tests

This project includes two unit tests that demonstrate key aspects of the server's functionality, particularly focusing on graceful shutdowns:

### 1. `TestGracefulShutdownSuccess`

This test simulates a successful graceful shutdown. It sets up an HTTP POST request with a `sleep` value of 15 seconds and then triggers a graceful shutdown after a short delay. The test verifies that the server shuts down gracefully within the allotted time. The test code uses the standard `testing` package to capture log output and ensure no errors occurred during the shutdown process.

### 2. `TestGracefulShutdownTimeout`

This test simulates a scenario where the server fails to shut down gracefully due to a very short timeout. Similar to the first test, it sends an HTTP POST request with a `sleep` value, but the timeout for the shutdown process is set to just 1 millisecond. This is intended to force a timeout and check that the server logs an appropriate message about not being able to shut down gracefully.

### Simulating an Interrupt Signal

An uncommon but important aspect of these tests is the simulation of an operating system interrupt signal to trigger the graceful shutdown. This is achieved using the `os` package in conjunction with the `syscall` package to send a `SIGINT` signal (the same signal sent by pressing `Ctrl+C` in a terminal). This approach showcases the flexibility of Go's `os` package and how you can simulate real-world scenarios in your unit tests.

```go
// Simulate an OS interrupt after a short delay
go func() {
    time.Sleep(500 * time.Millisecond)
    p, _ := os.FindProcess(os.Getpid())
    p.Signal(syscall.SIGINT)
}()
```

This technique is particularly useful for testing how your server responds to typical shutdown signals in a controlled environment, ensuring that your application behaves correctly in production scenarios.



## License

This demonstration code and project files is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for more details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you have suggestions or improvements.

## Author

[Ashwin Rao](https://github.com/yourusername)