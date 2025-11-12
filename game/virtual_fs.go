package game

import (
	"encoding/base64"
	"sort"
	"strings"
)

type VirtualFileSystem struct {
	files map[string]string
}

func NewVirtualFS() *VirtualFileSystem {
	return &VirtualFileSystem{
		files: make(map[string]string),
	}
}

func (vfs *VirtualFileSystem) ListFiles(path string) string {
	if path != "." && path != "" {
		return "ls: " + path + ": No such directory"
	}

	var result strings.Builder
	for filename := range vfs.files {
		// Skip hidden files in regular ls
		if !strings.HasPrefix(filename, ".") {
			result.WriteString(filename)
			result.WriteString("\n")
		}
	}
	return strings.TrimSpace(result.String())
}

// ListAllFiles shows ALL files including hidden ones (for ls -a)
func (vfs *VirtualFileSystem) ListAllFiles(path string) string {
	if path != "." && path != "" {
		return "ls: " + path + ": No such directory"
	}

	var files []string
	for filename := range vfs.files {
		files = append(files, filename)
	}

	// Sort files for consistent output
	sort.Strings(files)

	var result strings.Builder
	for _, filename := range files {
		result.WriteString(filename)
		result.WriteString("\n")
	}
	return strings.TrimSpace(result.String())
}

func (vfs *VirtualFileSystem) ReadFile(filename string) (string, bool) {
	content, exists := vfs.files[filename]
	return content, exists
}

func (vfs *VirtualFileSystem) WriteFile(filename, content string) {
	vfs.files[filename] = content
}

func (vfs *VirtualFileSystem) FindFiles(args []string) string {
	if len(args) == 0 {
		return "find: missing search pattern"
	}

	pattern := args[0]
	var results []string

	for filename := range vfs.files {
		if strings.Contains(filename, pattern) {
			results = append(results, filename)
		}
	}

	if len(results) == 0 {
		return "No files found matching: " + pattern
	}

	return strings.Join(results, "\n")
}

// --- NEW METHODS FOR CHALLENGING LEVELS ---

// GrepFile searches for pattern in file content (for level 5)
func (vfs *VirtualFileSystem) GrepFile(pattern, filename string) string {
	content, exists := vfs.files[filename]
	if !exists {
		return "grep: " + filename + ": No such file or directory"
	}

	lines := strings.Split(content, "\n")
	var matches []string

	for _, line := range lines {
		if strings.Contains(line, pattern) {
			matches = append(matches, line)
		}
	}

	if len(matches) == 0 {
		return "" // No output if no matches (standard grep behavior)
	}

	return strings.Join(matches, "\n")
}

// StringsCommand extracts readable strings from "binary" files (for level 6)
func (vfs *VirtualFileSystem) StringsCommand(filename string) string {
	content, exists := vfs.files[filename]
	if !exists {
		return "strings: '" + filename + "': No such file"
	}

	// Simulate extracting readable strings from binary data
	// In real binary, we'd filter only readable ASCII, but here we'll just return content
	return content
}

// Base64Decode decodes base64 encoded content (for level 9)
func (vfs *VirtualFileSystem) Base64Decode(filename string) string {
	encoded, exists := vfs.files[filename]
	if !exists {
		return "base64: " + filename + ": No such file or directory"
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "base64: invalid input"
	}

	return string(decoded)
}
