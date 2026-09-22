package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
)

type Event struct {
	req_type string
	data     string
	user     string
	conn     net.Conn
}

type Hub struct {
	events  chan Event
	clients map[string]net.Conn
	done    chan struct{}
}

func (h *Hub) run() {
	for {
		select {
		case event := <-h.events:

			switch event.req_type {

			case "JOIN":
				if _, exists := h.clients[event.user]; exists {
					fmt.Fprintln(event.conn, "Username already taken.")
					event.conn.Close()
					continue
				}

				h.clients[event.user] = event.conn

				fmt.Println(event.user, "joined the chat.")

				h.broadcast(
					event.user+" joined the chat.",
					event.user,
				)

				fmt.Fprintln(event.conn, "Welcome to ChatHub, "+event.user+"!")

			case "MESSAGE":
				fmt.Printf("%s: %s\n", event.user, event.data)

				h.broadcast(
					event.user+": "+event.data,
					event.user,
				)

			case "LEAVE":
				if _, exists := h.clients[event.user]; exists {
					delete(h.clients, event.user)

					fmt.Println(event.user, "left the chat.")

					h.broadcast(
						event.user+" left the chat.",
						event.user,
					)
				}
			}

		case <-h.done:
			fmt.Println("Hub shutting down...")

			for user, conn := range h.clients {
				fmt.Fprintln(conn, "Server shutting down.")
				conn.Close()
				delete(h.clients, user)
			}

			return
		}
	}
}

func (h *Hub) broadcast(message string, sender string) {
	for user, conn := range h.clients {
		if user == sender {
			continue
		}

		_, err := fmt.Fprintln(conn, message)
		if err != nil {
			fmt.Println("broadcast error:", err)
		}
	}
}

func handleConnection(conn net.Conn, h *Hub, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	reader := bufio.NewReader(conn)

	fmt.Fprintln(conn, "Enter username:")

	username, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	username = strings.TrimSpace(username)

	if username == "" {
		fmt.Fprintln(conn, "Username cannot be empty.")
		return
	}

	h.events <- Event{
		req_type: "JOIN",
		user:     username,
		conn:     conn,
	}

	fmt.Fprintln(conn, "You can start chatting.")

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			select {
			case h.events <- Event{
				req_type: "LEAVE",
				user:     username,
				conn:     conn,
			}:
			case <-h.done:
			}

			return
		}

		message = strings.TrimSpace(message)

		if message == "" {
			continue
		}

		if message == "/quit" {
			select {
			case h.events <- Event{
				req_type: "LEAVE",
				user:     username,
				conn:     conn,
			}:
			case <-h.done:
			}

			return
		}

		select {
		case h.events <- Event{
			req_type: "MESSAGE",
			data:     message,
			user:     username,
			conn:     conn,
		}:
		case <-h.done:
			return
		}
	}
}

func main() {
	h := Hub{
		events:  make(chan Event, 100),
		clients: make(map[string]net.Conn),
		done:    make(chan struct{}),
	}

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("ChatHub running on :9000")

	go h.run()

	var wg sync.WaitGroup

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan

		fmt.Println("\nShutting down...")

		close(h.done)
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()

		if err != nil {
			select {
			case <-h.done:
				wg.Wait()
				fmt.Println("Server stopped.")
				return
			default:
				fmt.Println("accept error:", err)
				continue
			}
		}

		wg.Add(1)

		go handleConnection(conn, &h, &wg)
	}
}
