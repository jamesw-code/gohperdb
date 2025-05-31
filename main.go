package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Hello, Gopher!")
	port := ":9090"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("GopherDb listening on %s...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	fmt.Printf("Client connected: %s\n", addr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("[%s] Received: %s\n", addr, text)

		// Echo back the message
		_, _ = conn.Write([]byte("Echo: " + text + "\n"))
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[%s] Connection error: %v\n", addr, err)
	} else {
		fmt.Printf("Client disconnected: %s\n", addr)
	}
}
