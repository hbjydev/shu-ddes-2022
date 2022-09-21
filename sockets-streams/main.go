package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"net"
    "sockets-streams/http"
	"strings"
)

var (
    // Allow the listen address to be configurable on the command line
    // by using -listen=:8080
    listenAddr = flag.String("listen", ":8000", "the address to listen on")
)

func main() {
    flag.Parse()
	log.Printf(`starting server on %s`, *listenAddr)

	// Create a socket on port 8000, bound to all ifaces
	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("failed to bind to %s: %s", *listenAddr, err.Error())
	}
    defer lis.Close()

    // Constantly listen for incoming connections
    for {
        // Once a new connection has been established, accept the connection.
        conn, err := lis.Accept()
        if err != nil {
            // Kill the server if accepting the socket connection fails.
            log.Fatalf("failed to accept connection: %s", err.Error())
        }

        // Handle the TCP connection in a parallel process.
        //
        // This example uses the Go feature called Goroutines, Go's lightweight
        // threading system.
        //
        // It spawns a new thread to run the function handleRequest, given the
        // net.Conn instance created by lis.Accept()
        //
        // - Ref: https://go.dev/ref/spec#Go_statements
        // - Guide: https://go.dev/tour/concurrency
        go handleRequest(conn)
    }
}

// handleRequest handles a raw TCP connection and deals with the flow of
// streaming input and output across the opened socket.
func handleRequest(conn net.Conn) {
    log.Printf("new connection from %v", conn.RemoteAddr())

    // Close the connection once the routine finishes
    defer conn.Close()

    lines := []string{}

    for {
        // Read the data into a string, delimited by a newline
        line, err := bufio.NewReader(conn).ReadString('\n')

        // If the error reading is an early disconnect by the client,
        // just break out of the read loop
        if err != nil {
            if err.Error() == "EOF" {
                conn.Close()
                return
            }

            // Otherwise, print the error
            log.Printf("error processing connection data buffer: %s", err.Error())
            break
        }

        // If the current request is the stop code, end the TCP session
        if line == "\n" {
            break
        }

        lines = append(lines, strings.TrimSpace(line))
    }

    req, err := http.Parse(lines)

    if err != nil {
        response := ApiResponse[any]{
            Message: "failed to parse http query.",
            Code: 400,
            Data: nil,
        }

        responseBytes, err := json.Marshal(response)
        if err != nil {
            log.Printf("could not marshal api response: %s", err.Error())
            return
        }

        result := string(responseBytes) + "\n"
        conn.Write([]byte(result))       
    }

    response := ApiResponse[any]{
        Message: "Hello, world.",
        Code: 1,
        Data: req,
    }

    responseBytes, err := json.Marshal(response)
    if err != nil {
        log.Printf("could not marshal api response: %s", err.Error())
        return
    }

    result := string(responseBytes) + "\n"
    conn.Write([]byte(result))
}

type ApiResponse [T any] struct{
    Message string `json:"message"`
    Code int `json:"code"`
    Data T `json:"data"`
}

