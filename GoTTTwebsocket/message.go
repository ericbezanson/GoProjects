package message

// Message Struct for data being sent over websocket
// Type: Describes the type of message, such as a chat message or a move in the game.
// Text: The content of the message (e.g., "User X has joined the game").
// Sender: An optional field for the sender's name.
// UserName: An optional field for the user's unique name (e.g., user-1).
// Symbol: An optional field for the player's symbol, either "X" or "O".
// Position: A pointer to an integer, representing the position on the Tic-Tac-Toe board (optional and can be nil).
type Message struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	Sender   string `json:"sender,omitempty"`
	UserName string `json:"userName,omitempty"`
	Symbol   string `json:"symbol,omitempty"`
	Position int    `json:"position"` // Allow for zero int value
}
