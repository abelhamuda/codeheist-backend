package websocket

import (
	"log"
	"net/http"
	"strings"

	"codeheist/game"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type WebSocketHandler struct {
	engine *game.GameEngine
}

type WSMessage struct {
	Type      string `json:"type"`
	Content   string `json:"content,omitempty"`
	Data      string `json:"data,omitempty"`
	Command   string `json:"command,omitempty"`
	Level     int    `json:"level,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}

func NewHandler(engine *game.GameEngine) *WebSocketHandler {
	return &WebSocketHandler{
		engine: engine,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Create new session
	ip := strings.Split(c.Request.RemoteAddr, ":")[0] // Get IP without port
	session := h.engine.CreateSession(ip)

	log.Printf("🔗 New WebSocket connection from %s, session: %s", ip, session.ID)

	// Send session created message
	welcomeMsg := WSMessage{
		Type:      "session_created",
		Content:   "\r\n\x1b[32m● WELCOME TO CODEHEIST\x1b[0m\r\n" + getLevelWelcomeMessage(h.engine, session.CurrentLevel) + "\r\n",
		SessionID: session.ID,
	}
	conn.WriteJSON(welcomeMsg)

	// Send initial prompt
	promptMsg := WSMessage{
		Type:    "prompt",
		Content: "$ ",
	}
	conn.WriteJSON(promptMsg)

	// Handle messages from client
	for {
		var msg WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		log.Printf("📨 Received message type: %s", msg.Type)

		switch msg.Type {
		case "command":
			h.handleCommand(conn, session.ID, msg.Command)
		case "command_input":
			// Handle direct input from terminal
			h.handleCommandInput(conn, session.ID, msg.Data)
		default:
			log.Printf("❌ Unknown message type: %s", msg.Type)
		}
	}
}

// Handle command input (character by character)
func (h *WebSocketHandler) handleCommandInput(conn *websocket.Conn, sessionID string, input string) {
	// For now, we'll handle complete commands only
	// This can be extended for real-time input handling
	if input == "\r" || input == "\n" { // Enter key
		session, exists := h.engine.GetSession(sessionID)
		if exists && session.CurrentInput != "" {
			h.handleCommand(conn, sessionID, session.CurrentInput)
			session.CurrentInput = ""
		} else {
			// Send empty prompt
			promptMsg := WSMessage{
				Type:    "prompt",
				Content: "$ ",
			}
			conn.WriteJSON(promptMsg)
		}
	} else if input == "\x7f" { // Backspace
		session, exists := h.engine.GetSession(sessionID)
		if exists && len(session.CurrentInput) > 0 {
			session.CurrentInput = session.CurrentInput[:len(session.CurrentInput)-1]
		}
	} else {
		// Append to current input
		session, exists := h.engine.GetSession(sessionID)
		if exists {
			session.CurrentInput += input
		}
	}
}

// Helper function to get level welcome message
func getLevelWelcomeMessage(engine *game.GameEngine, level int) string {
	if level < len(engine.Levels) {
		return engine.Levels[level].WelcomeMsg
	}
	return "Welcome to CodeHeist! Your mission awaits..."
}

func (h *WebSocketHandler) handleCommand(conn *websocket.Conn, sessionID string, command string) {
	log.Printf("🔧 Executing command: '%s' for session: %s", command, sessionID)

	// Execute command
	response := h.engine.ExecuteCommand(sessionID, command)

	// Send command output
	if response.Output != "" {
		outputMsg := WSMessage{
			Type:    "output",
			Content: response.Output + "\r\n",
		}
		conn.WriteJSON(outputMsg)
	}

	// Handle level completion
	if response.LevelCompleted {
		levelUpMsg := WSMessage{
			Type:  "level_up",
			Level: response.NewLevel,
		}
		conn.WriteJSON(levelUpMsg)

		// Send welcome message for new level
		session, exists := h.engine.GetSession(sessionID)
		if exists {
			welcomeMsg := WSMessage{
				Type:    "output",
				Content: "\r\n\x1b[36m" + getLevelWelcomeMessage(h.engine, session.CurrentLevel) + "\x1b[0m\r\n",
			}
			conn.WriteJSON(welcomeMsg)
		}
	}

	// Always send new prompt after command execution
	promptMsg := WSMessage{
		Type:    "prompt",
		Content: "$ ",
	}
	conn.WriteJSON(promptMsg)
}
