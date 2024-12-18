package game

import (
	"fmt"

	"github.com/ericbezanson/GoProjects/tree/main/GoTTTwebsocket/message"
	"golang.org/x/net/websocket"
)

// Tic-Tac-Toe board representation (fixed-size array)
var board [9]string

// Track the current player (default "X")
var currentPlayer = "X"

// Game state flags
var gameStarted bool

// args
// ws *websocket.Conn: The WebSocket connection of the player making the move. (pointer)
// position *int: A pointer to the board position where the player wants to place their symbol (accept 0 value).
// symbol string: The player’s symbol ("X" or "O").
func handleMove(ws *websocket.Conn, position int, sender string, symbol string) {
	// Validate the move
	if position < 0 || position > 8 {
		fmt.Println("Invalid move: Position is out of bounds")
		return
	}
	if currentPlayer != symbol {
		fmt.Println("Invalid move: Not your turn")
		return
	}
	if board[position] != "" {
		fmt.Println("Invalid move: Cell already occupied")
		return
	}

	// Update the game board / place symbol in clicked tile
	board[position] = symbol

	// Broadcast the move to all clients
	sendMessageToAll(message.Message{
		Type:     "move",
		Position: position,
		Text:     symbol,
	})

	// Check if the current move resulted in a win
	if winPattern := checkWin(symbol); len(winPattern) > 0 {
		// Announce the winner
		sendMessageToAll(message.Message{
			Type:     "gameOver",
			Text:     fmt.Sprintf("User-%s Wins!", symbol),
			Symbol:   symbol,
			Position: -1, // Unused
		})

		// Reset the game
		resetGame()
		return
	}

	// If no win, check for a draw
	if checkStalemate() {
		sendMessageToAll(message.Message{
			Type: "gameOver",
			Text: "It's a draw!",
		})
		resetGame()
		return
	}

	// Switch turns
	switchTurn()

	// Notify players of the turn change
	sendMessageToAll(message.Message{Type: "updateTurn", Text: currentPlayer})
}

// Check if the given symbol has won
func checkWin(symbol string) [][3]int {
	// all possible win patterns in tic tac toe
	// slice of arrays [3]int
	winPatterns := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}

	// iterate over all win patterns and check if the player has won
	var winningPatterns [][3]int
	for _, pattern := range winPatterns {
		if board[pattern[0]] == symbol && board[pattern[1]] == symbol && board[pattern[2]] == symbol {
			winningPatterns = append(winningPatterns, pattern)
		}
	}
	return winningPatterns
}

func checkStalemate() bool {
	for _, cell := range board {
		if cell == "" {
			return false // There's still an empty cell
		}
	}
	return true
}

func resetGame() {
	// Reset the board and game state
	board = [9]string{"", "", "", "", "", "", "", "", ""}
	gameStarted = false
	userCount = 0
	currentPlayer = "X"
}

func switchTurn() {
	if currentPlayer == "X" {
		currentPlayer = "O"
	} else {
		currentPlayer = "X"
	}
}

// Function to send the initial board state to a client
func SendBoardState(ws *websocket.Conn) {
	sendMessage(ws, message.Message{
		Type: "updateBoard",
		Text: fmt.Sprintf("%v", board),
	})
}

// Function to send a message to a specific client
func sendMessage(ws *websocket.Conn, msg message.Message) {
	websocket.JSON.Send(ws, msg)
}
