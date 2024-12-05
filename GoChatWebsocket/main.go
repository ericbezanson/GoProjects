package main

import (
	"fmt"      // Provides functions for formatted I/O,
	"net/http" // Provides HTTP client and server functionality.

	"github.com/gorilla/websocket" // The websocket package from the Gorilla toolkit is used for handling WebSocket connections,
)

// --- GLOBALS ---

// instance of websocket.Upgrader, used to upgrade an HTTP connection to a Websocket one in handleConnections()
var websocketUpgrader = websocket.Upgrader{
	// return true on check origin to allow requests from any origin, for demo purposes (not very secure)
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

var clients = make(map[*websocket.Conn]bool) // Connected clients

var broadcast = make(chan string) // Broadcast channel

// --- FUNCTIONS ---

// app entry
func main() {
	// Sets up a route (/) that serves static files. http.FileServer(http.Dir("./")) serves the contents of the current directory (index.html)
	http.Handle("/", http.FileServer(http.Dir("./")))

	// sets up the /ws endpoint to handle WebSocket connections. When a WebSocket request is made to this endpoint, it triggers the handleConnections function to manage the connection.
	http.HandleFunc("/ws", handleConnections)

	// Starts the handleMessages function as a goroutine. This function listens for messages on the broadcast channel and sends them to all connected clients.
	go handleMessages()

	// Starts the HTTP server on port 8080. The server will handle incoming requests (including WebSocket and static file requests) on this port.
	fmt.Println("Chat server started on :8080")
	err := http.ListenAndServe(":8080", nil)

	// Error Handling
	if err != nil {
		fmt.Println("Server Error:", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {

	// use websocket.Upgrader.Upgrade to upgrade from an http connection to a websocket connection
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}

	// ensure connection is closed when the function exits (client D/C or error etc)
	defer conn.Close()

	// Log the new connection
	fmt.Println("New WebSocket connection established")

	// Adds the WebSocket connection to the clients map to track active connections.
	clients[conn] = true

	// listen for messages inside the loop
	for {
		// reads a message from the WebSocket connection. The msgBytes variable holds the byte slice (in ASCII) of the message.
		// conn.ReadMessage:
		_, msgBytes, err := conn.ReadMessage()

		// error handling
		// if error, remove connection from clients list / exit loop
		if err != nil {
			fmt.Println("Error reading message:", err)
			delete(clients, conn)
			break
		}

		// Convert the byte slice (ASC II chars) into a readable string
		msg := string(msgBytes)

		// message is then sent to the broadcast channel
		// <- :
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// constantly listen for messages on the broadcast channel, when found assign to msg variable
		msg := <-broadcast

		// loop through all connected clients
		for client := range clients {
			// write the message to each individual client in JSON format via WriteJSON
			err := client.WriteJSON(msg)

			// if error log error, close client connection and remove connection from client list
			if err != nil {
				fmt.Println("Error writing to client:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
