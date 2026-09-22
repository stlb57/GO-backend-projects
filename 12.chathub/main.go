package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

type Event struct {
	req_type string
	data     string
	user     string
	conn     net.Conn
}

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

type Hub struct {
	events   chan Event
	clients  map[string]net.Conn
	messages []Message
	done     chan struct{}
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

				h.broadcast(event.user+" joined the chat.", event.user)

				fmt.Fprintln(
					event.conn,
					"Welcome to ChatHub, "+event.user+"!",
				)

			case "MESSAGE":
				fmt.Printf("%s: %s\n", event.user, event.data)

				h.messages = append(h.messages, Message{
					User:    event.user,
					Message: event.data,
				})

				h.broadcast(
					event.user+": "+event.data,
					event.user,
				)

			case "BROADCAST":
				fmt.Println("SERVER:", event.data)

				h.broadcast(
					"SERVER: "+event.data,
					"",
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

func handleConnection(
	conn net.Conn,
	h *Hub,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	defer conn.Close()

	fmt.Fprintln(conn, "Enter username:")

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	username := strings.TrimSpace(string(buffer[:n]))

	if username == "" {
		fmt.Fprintln(conn, "Username cannot be empty.")
		return
	}

	select {
	case h.events <- Event{
		req_type: "JOIN",
		user:     username,
		conn:     conn,
	}:
	case <-h.done:
		return
	}

	fmt.Fprintln(conn, "You can start chatting.")

	for {

		n, err := conn.Read(buffer)

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

		message := strings.TrimSpace(
			string(buffer[:n]),
		)

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

// ---------------- HTTP ----------------

type HTTPHandler struct {
	hub *Hub
}

func (h HTTPHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.URL.Path {

	case "/":
		h.home(w, r)

	case "/health":
		h.health(w, r)

	case "/users":
		h.users(w, r)

	case "/messages":
		h.messagesHandler(w, r)

	case "/broadcast":
		h.broadcastHandler(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (h HTTPHandler) home(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set(
		"Content-Type",
		"text/plain",
	)

	fmt.Fprintln(w, "Welcome to ChatHub")
}

func (h HTTPHandler) health(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set(
		"Content-Type",
		"text/plain",
	)

	fmt.Fprintln(w, "OK")
}

func (h HTTPHandler) users(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	users := make([]string, 0, len(h.hub.clients))

	for user := range h.hub.clients {
		users = append(users, user)
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string][]string{
			"users": users,
		},
	)
}

func (h HTTPHandler) messagesHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string][]Message{
			"messages": h.hub.messages,
		},
	)
}

func (h HTTPHandler) broadcastHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Message string `json:"message"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(
			w,
			"invalid JSON",
			http.StatusBadRequest,
		)

		return
	}

	if strings.TrimSpace(request.Message) == "" {
		http.Error(
			w,
			"message cannot be empty",
			http.StatusBadRequest,
		)

		return
	}

	select {
	case h.hub.events <- Event{
		req_type: "BROADCAST",
		data:     request.Message,
	}:
	case <-h.hub.done:
		http.Error(
			w,
			"server shutting down",
			http.StatusServiceUnavailable,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"status": "broadcast sent",
		},
	)
}

// ---------------- MAIN ----------------

func main() {

	hub := Hub{
		events:   make(chan Event, 100),
		clients:  make(map[string]net.Conn),
		messages: make([]Message, 0),
		done:     make(chan struct{}),
	}

	// TCP SERVER

	listener, err := net.Listen(
		"tcp",
		":9000",
	)

	if err != nil {
		fmt.Println("TCP listen error:", err)
		return
	}

	defer listener.Close()

	fmt.Println("ChatHub TCP server running on :9000")

	go hub.run()
	mux := http.NewServeMux()

	httpHandler := HTTPHandler{
		hub: &hub,
	}

	mux.Handle("/", httpHandler)

	httpServer := http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		fmt.Println("ChatHub HTTP server running on :8080")

		err := httpServer.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			fmt.Println("HTTP server error:", err)
		}
	}()

	// SHUTDOWN

	var wg sync.WaitGroup

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		os.Interrupt,
	)

	<-signalChan

	fmt.Println("\nShutting down ChatHub...")

	close(hub.done)

	listener.Close()

	httpServer.Close()

	wg.Wait()

	fmt.Println("ChatHub stopped.")
}
