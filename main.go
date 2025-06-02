package main

import (
	"bufio"
	"fmt"
	"gopherdb/store"
	"net"
	"strconv"
	"strings"
	"time"
)

var backend = store.NewInMemoryStore()

func main() {
	fmt.Println("Hello, Gopher!")

	port := ":5321"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("GopherDb listening on %s...\n", port)

	go func() {
		for {
			backend.CleanupExpired()
			time.Sleep(1 * time.Minute)
		}
	}()

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
		parts := strings.Fields(text)
		fmt.Printf("[%s] Received: %s\n", addr, text)

		if len(parts) == 0 {
			return
		}

		switch strings.ToUpper(parts[0]) {
		case "SET":
			if len(parts) != 3 {
				conn.Write([]byte("ERR: usage SET key value\n"))
				return
			}
			backend.Set(parts[1], parts[2])
			conn.Write([]byte("OK\n"))

		case "GET":
			if len(parts) != 2 {
				conn.Write([]byte("ERR: usage GET key\n"))
				return
			}
			val, ok := backend.Get(parts[1])
			if !ok {
				conn.Write([]byte("NULL\n"))
			} else {
				conn.Write([]byte(val + "\n"))
			}

		case "DELETE":
			if len(parts) != 2 {
				conn.Write([]byte("ERR: usage DELETE key\n"))
				return
			}
			backend.Delete(parts[1])
			conn.Write([]byte("OK\n"))

		case "EXISTS":
			if len(parts) != 2 {
				conn.Write([]byte("ERR: usage EXISTS key\n"))
				return
			}
			_, ok := backend.Get(parts[1])
			if ok {
				conn.Write([]byte("1\n"))
			} else {
				conn.Write([]byte("0\n"))
			}
		case "SETEX":
			if len(parts) != 4 {
				conn.Write([]byte("ERR: usage SETEX key ttl value\n"))
				return
			}

			ttl, err := strconv.Atoi(parts[2])
			if err != nil {
				conn.Write([]byte("ERR: TTL must be an integer\n"))
				return
			}

			backend.SetEX(parts[1], parts[3], ttl)
			conn.Write([]byte("OK\n"))

		case "PING":
			conn.Write([]byte("PONG\n"))

		default:
			_, _ = conn.Write([]byte("ERR: unknown command\n"))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[%s] Connection error: %v\n", addr, err)
	} else {
		fmt.Printf("Client disconnected: %s\n", addr)
	}
}
