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
		User:         "codeheist0",
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

	// CHECK: Jika level tidak ada, berarti game completed!
	if s.CurrentLevel >= len(e.Levels) {
		return &CommandResponse{
			Output: "\n🎉 CONGRATULATIONS! You've completed all levels!\n" +
				"🏆 You are now a CodeHeist Master!\n\n" +
				"Thank you for playing! 🚀",
		}
	}

	level := e.Levels[s.CurrentLevel]
	output, levelCompleted := e.processCommand(command, s, level)

	if levelCompleted {
		oldLevel := s.CurrentLevel
		s.CurrentLevel++

		// CHECK: Jika masih ada level berikutnya, initialize filesystem
		if s.CurrentLevel < len(e.Levels) {
			e.initializeLevelFilesystem(s, s.CurrentLevel)
		}

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

// parseCommandWithQuotes handles command parsing with proper quote support
func parseCommandWithQuotes(input string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(input); i++ {
		c := input[i]

		switch {
		case c == '"' || c == '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = c
			} else if c == quoteChar {
				inQuotes = false
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			} else {
				current.WriteByte(c)
			}
		case c == ' ' && !inQuotes:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(c)
		}
	}

	// Handle last argument
	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

func (e *GameEngine) processCommand(cmd string, session *Session, level *Level) (string, bool) {
	cmd = strings.TrimSpace(cmd)

	// Use the new quotes-aware parser instead of strings.Fields
	parts := parseCommandWithQuotes(cmd)

	if len(parts) == 0 {
		return "", false
	}

	command := parts[0]
	args := parts[1:]

	// Variables to store command output and completion status
	var output string
	completed := false

	switch command {
	case "ls":
		// Handle ls flags like -a
		if len(args) == 0 {
			output = session.VirtualFS.ListFiles(".")
		} else if len(args) == 1 && args[0] == "-a" {
			output = session.VirtualFS.ListAllFiles(".")
		} else {
			output = session.VirtualFS.ListFiles(args[0])
		}

	case "cat":
		if len(args) == 0 {
			output = "cat: missing filename"
			break
		}

		filename := args[0]
		content, exists := session.VirtualFS.ReadFile(filename)
		if !exists {
			output = "cat: " + filename + ": No such file or directory"
			break
		}
		output = content
		completed = (content == level.Solution)

	case "cd":
		// Simple cd implementation
		if len(args) == 0 || args[0] == "~" {
			output = ""
		} else {
			output = "cd: " + args[0] + ": No such directory"
		}

	case "find":
		output = session.VirtualFS.FindFiles(args)

	case "hint":
		output = "💡 Hint: " + level.Hint

	case "status":
		output = e.getStatus(session)

	// --- NEW COMMANDS FOR CHALLENGING LEVELS ---
	case "chmod":
		if len(args) < 2 {
			output = "chmod: missing operand"
			break
		}
		// Simulate permission change for level 4
		if args[1] == "secret.txt" {
			output = "Permissions changed for secret.txt"
			// After chmod, the file becomes readable with cat
		} else {
			output = "chmod: cannot access '" + args[1] + "': No such file or directory"
		}

	case "grep":
		if len(args) < 2 {
			output = "grep: missing pattern or filename"
			break
		}
		output = session.VirtualFS.GrepFile(args[0], args[1])
		// Check if grep output contains the EXACT solution
		if strings.Contains(output, level.Solution) {
			// Extract just the solution line for cleaner completion
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if strings.Contains(line, level.Solution) {
					output = strings.TrimSpace(line)
					completed = true
					break
				}
			}
		}

	case "strings":
		if len(args) == 0 {
			output = "strings: missing filename"
			break
		}
		output = session.VirtualFS.StringsCommand(args[0])
		// For strings command, check if output contains the exact solution
		if strings.Contains(output, level.Solution) {
			// Extract just the solution part
			start := strings.Index(output, level.Solution)
			if start != -1 {
				end := start + len(level.Solution)
				output = output[start:end]
				completed = true
			}
		}

	case "echo":
		if len(args) == 0 {
			output = ""
			break
		}
		// Handle environment variables for level 8
		if len(args) == 1 && args[0] == "$SECRET_KEY" {
			output = "bandit9{EnvVariableMaster}"
			completed = true
		} else {
			output = strings.Join(args, " ")
		}

	case "base64":
		if len(args) < 2 {
			output = "base64: missing option or filename"
			break
		}
		if args[0] == "-d" {
			output = session.VirtualFS.Base64Decode(args[1])
			// Check if decoded content matches solution exactly
			completed = (strings.TrimSpace(output) == level.Solution)
		} else {
			output = "base64: invalid option -- '" + args[0] + "'"
		}

	default:
		output = "command not found: " + command
	}

	// FINAL CHECK: Only complete if we haven't already marked as completed
	// AND output exactly matches the solution (to prevent false positives)
	if !completed {
		cleanOutput := strings.TrimSpace(output)
		if cleanOutput == level.Solution {
			completed = true
		}
	}

	if completed {
		return output + "\n\n" + levelCompletionMessage(level.ID), true
	}

	return output, false
}

func (e *GameEngine) showHelp() *CommandResponse {
	helpText := `Available commands:
  ls [dir]        - List directory contents
  ls -a          - List all files including hidden
  cat <file>      - Display file contents  
  cd [dir]        - Change directory
  find [pattern]  - Find files
  grep <pattern> <file> - Search for text in files
  strings <file>  - Extract text from binary files
  chmod <mode> <file> - Change file permissions
  echo <text>     - Display text or variables
  base64 -d <file> - Decode base64 encoded file
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
	// Safety check - pastikan level exists
	if level < 0 || level >= len(e.Levels) {
		log.Printf("⚠️ Attempted to initialize invalid level: %d", level)
		return
	}

	levelConfig := e.Levels[level]
	session.VirtualFS = NewVirtualFS()
	session.User = fmt.Sprintf("codeheist%d", level)

	// Populate filesystem for this level
	for filename, content := range levelConfig.Filesystem {
		session.VirtualFS.WriteFile(filename, content.(string))
	}

	log.Printf("🔄 Initialized filesystem for level %d", level)
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
