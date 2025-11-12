package game

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type GameEngine struct {
	Sessions sync.Map
	Levels   map[int]*Level
}

type Session struct {
	ID           string
	CurrentLevel int
	VirtualFS    *VirtualFileSystem
	User         string
	CreatedAt    time.Time
	IPAddress    string
	LastActivity time.Time
	CurrentInput string
	mu           sync.Mutex // Add mutex for CurrentInput safety
}

type Level struct {
	ID          int
	Title       string
	Description string
	Filesystem  map[string]interface{}
	Solution    string
	Hint        string
	WelcomeMsg  string
}

type CommandResponse struct {
	Output         string
	LevelCompleted bool
	NewLevel       int
}

func NewEngine() *GameEngine {
	engine := &GameEngine{
		Levels: initializeLevels(),
	}
	return engine
}

func (e *GameEngine) CreateSession(ip string) *Session {
	session := &Session{
		ID:           uuid.New().String(),
		CurrentLevel: 0,
		VirtualFS:    NewVirtualFS(),
		User:         "bandit0",
		CreatedAt:    time.Now(),
		IPAddress:    ip,
		LastActivity: time.Now(),
	}

	e.initializeLevelFilesystem(session, 0)
	e.Sessions.Store(session.ID, session)

	log.Printf("🆕 New session created: %s for %s", session.ID, ip)
	return session
}

func (e *GameEngine) GetSession(sessionID string) (*Session, bool) {
	session, exists := e.Sessions.Load(sessionID)
	if !exists {
		return nil, false
	}
	return session.(*Session), true
}

func (e *GameEngine) ExecuteCommand(sessionID, command string) *CommandResponse {
	session, exists := e.Sessions.Load(sessionID)
	if !exists {
		return &CommandResponse{Output: "Session not found"}
	}

	s := session.(*Session)
	s.LastActivity = time.Now()

	// Handle special commands
	switch strings.TrimSpace(command) {
	case "clear":
		return &CommandResponse{Output: "\033[H\033[2J"}
	case "help":
		return e.showHelp()
	case "levels":
		return e.showLevels()
	case "whoami":
		return &CommandResponse{Output: s.User}
	case "pwd":
		return &CommandResponse{Output: "/home/" + s.User}
	}

	level := e.Levels[s.CurrentLevel]
	output, levelCompleted := e.processCommand(command, s, level)

	if levelCompleted {
		oldLevel := s.CurrentLevel
		s.CurrentLevel++
		e.initializeLevelFilesystem(s, s.CurrentLevel)

		log.Printf("🎉 Session %s completed level %d", sessionID, oldLevel)

		return &CommandResponse{
			Output:         output,
			LevelCompleted: true,
			NewLevel:       s.CurrentLevel,
		}
	}

	return &CommandResponse{
		Output:         output,
		LevelCompleted: false,
	}
}

func (e *GameEngine) processCommand(cmd string, session *Session, level *Level) (string, bool) {
	cmd = strings.TrimSpace(cmd)
	parts := strings.Fields(cmd)

	if len(parts) == 0 {
		return "", false
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "ls":
		if len(args) == 0 {
			return session.VirtualFS.ListFiles("."), false
		}
		return session.VirtualFS.ListFiles(args[0]), false

	case "cat":
		if len(args) == 0 {
			return "cat: missing filename", false
		}
		content, exists := session.VirtualFS.ReadFile(args[0])
		if !exists {
			return "cat: " + args[0] + ": No such file or directory", false
		}

		// Check if this completes the level
		if content == level.Solution {
			return content + "\n\n" + levelCompletionMessage(level.ID), true
		}
		return content, false

	case "cd":
		// Simple cd implementation
		if len(args) == 0 || args[0] == "~" {
			return "", false // Change to home directory
		}
		return "cd: " + args[0] + ": No such directory", false

	case "find":
		return session.VirtualFS.FindFiles(args), false

	case "hint":
		return "💡 Hint: " + level.Hint, false

	case "status":
		return e.getStatus(session), false

	default:
		return "command not found: " + command, false
	}
}

func (e *GameEngine) showHelp() *CommandResponse {
	helpText := `Available commands:
  ls [dir]        - List directory contents
  cat <file>      - Display file contents  
  cd [dir]        - Change directory
  find [pattern]  - Find files
  pwd            - Print working directory
  whoami         - Show current user
  hint           - Get hint for current level
  status         - Show game status
  levels         - List all levels
  clear          - Clear terminal
  help           - Show this help message

Use these commands to find passwords and complete levels!`
	return &CommandResponse{Output: helpText}
}

func (e *GameEngine) showLevels() *CommandResponse {
	var levels []string
	for i := 0; i < len(e.Levels); i++ {
		level := e.Levels[i]
		levels = append(levels,
			fmt.Sprintf("Level %d: %s", level.ID, level.Title))
	}
	return &CommandResponse{Output: strings.Join(levels, "\n")}
}

func (e *GameEngine) getStatus(session *Session) string {
	level := e.Levels[session.CurrentLevel]
	return fmt.Sprintf(`
Current Level: %d
User: %s
Level Title: %s
Description: %s
Progress: %d/%d levels completed
	`,
		session.CurrentLevel,
		session.User,
		level.Title,
		level.Description,
		session.CurrentLevel,
		len(e.Levels)-1)
}

func (e *GameEngine) initializeLevelFilesystem(session *Session, level int) {
	levelConfig := e.Levels[level]
	session.VirtualFS = NewVirtualFS()
	// Fixed: Use fmt.Sprintf for proper string conversion
	session.User = fmt.Sprintf("bandit%d", level)

	// Populate filesystem for this level
	for filename, content := range levelConfig.Filesystem {
		session.VirtualFS.WriteFile(filename, content.(string))
	}
}

func (e *GameEngine) CleanupSessions() {
	ticker := time.NewTicker(30 * time.Minute)
	for range ticker.C {
		now := time.Now()
		e.Sessions.Range(func(key, value interface{}) bool {
			session := value.(*Session)
			if now.Sub(session.LastActivity) > 2*time.Hour {
				e.Sessions.Delete(key)
				log.Printf("🧹 Cleaned up expired session: %s", session.ID)
			}
			return true
		})
	}
}

func levelCompletionMessage(level int) string {
	messages := []string{
		"🎉 Excellent! Level completed!",
		"🚀 Great job! Moving to next challenge...",
		"💻 Hacking skills improving!",
		"🔓 Access granted to next level!",
		"⚡ Impressive terminal skills!",
	}

	if level < len(messages) {
		return messages[level]
	}
	return messages[len(messages)-1]
}
