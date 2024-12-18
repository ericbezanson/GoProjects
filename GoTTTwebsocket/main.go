package main

import (
	"fmt"
	"net/http"

	"github.com/ericbezanson/GoProjects/tree/main/GoTTTwebsocket/game"
	"github.com/ericbezanson/GoProjects/tree/main/GoTTTwebsocket/message"

	"golang.org/x/net/websocket" // switched from gorilla
)

// Globals

// Go map: data structure that acts as a collection of unordered key-value pairs
// use map here because we needed to store a key value pair (*websocket.Conn as a unique key) of a dynamic size
// NOTE: use the memory address as a unique key in the map
var clients = make(map[*websocket.Conn]string)

// same as clients, however this will store connection pointers for players who are not able to interact with game board
var spectators = make(map[*websocket.Conn]string)

// spectator count - used in spectator naming
var spectatorCount int

// A fixed-size array of strings representing the Tic-Tac-Toe board. Each element can be empty (""), "X", or "O".
// NOTE: more effecient to use a fixed sized array as TTT board will always be 3x3
var board = [9]string{"", "", "", "", "", "", "", "", ""}

// Track the number of users connected to the game.
var userCount int

// track if game is started, (needs to players)
var gameStarted bool

// keeps track of current player, default X for player one
var currentPlayer = "X"

func main() {
	// Sets up a handler to serve static files from current dir ./
	http.Handle("/", http.FileServer(http.Dir("./")))

	// sets up a WebSocket handler at the /ws path.
	http.Handle("/ws", websocket.Handler(handleConnections))

	game.SendBoardState(ws)
	// Starts the HTTP server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server Error:", err)
	}
}

func handleConnections(ws *websocket.Conn) {
	// defers the execution until the surrounding function returns.
	defer ws.Close()

	// username bucket, can be a user or spectator
	var userName string

	// Register user
	// if more than 2 users, spectator role assigned
	if userCount >= 2 {
		spectatorCount++
		userName = fmt.Sprintf("spectator-%d", spectatorCount)
		clients[ws] = userName

		// Notify spectator of status
		sendMessage(ws, message.Message{
			Type:     "lobbyFull",
			Text:     "The game lobby is full. You are now spectating.",
			UserName: userName,
		})

		// Broadcast spectator join message
		sendSystemMessage(fmt.Sprintf("%s has joined as a spectator.", userName))
	} else {
		// if not a spectator, assign user role
		userCount++
		userName = fmt.Sprintf("player-%d", userCount)
		clients[ws] = userName

		// Assign player symbol
		assignSymbol := "X"
		if userCount == 2 {
			assignSymbol = "O"
		}

		// Notify player of assignment
		sendMessage(ws, Message{
			Type:     "assignPlayer",
			UserName: userName,
			Symbol:   assignSymbol,
		})

		// Broadcast player join message
		sendSystemMessage(fmt.Sprintf("%s has joined the game.", userName))

		// Start the game when two players have joined
		if userCount == 2 && !gameStarted {
			gameStarted = true
			sendSystemMessage("Game has started! It's X's turn.")
			sendMessageToAll(message.Message{Type: "updateTurn", Text: "X"})
		}
	}

	// Send the initial board state to the new user
	sendMessage(ws, message.Message{
		Type: "updateBoard",
		Text: fmt.Sprintf("%v", board),
	})

	// Listen for messages
	for {
		var msg Message
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			fmt.Println("Connection closed:", err)
			delete(clients, ws)
			break
		}

		// Handle chat or move messages
		switch msg.Type {
		case "chat":
			sendMessageToAll(msg)
		case "move":
			handleMove(ws, msg.Position, msg.UserName, msg.Symbol)
		}
	}
}

func sendMessage(ws *websocket.Conn, msg message.Message) {
	websocket.JSON.Send(ws, msg)
}

func sendMessageToAll(msg message.Message) {
	fmt.Printf("Broadcasting message: %+v\n", msg)
	for client := range clients {
		sendMessage(client, msg)
	}
	for spectator := range spectators {
		sendMessage(spectator, msg)
	}
}

func sendSystemMessage(text string) {
	sendMessageToAll(Message{Type: "system", Text: text})
}

// TODO
// ------ MAJOR ------
// - Add start button / player ready
// - allow spectator to see board state if joining midgame
// - graceful shut down
// - update hardcoded localhost
// - isolate into new files (types)
// - add custom names
// - BUG: player 1 can interact with game board before game starts
// -----
// - broadcast state
// - keep gamestate checks on server side
// - multi games
// - lobby system
// - unit test, check win
// - table driven tests

// ------ MINOR ------
// highlight winning pattern
// chat with enter button
