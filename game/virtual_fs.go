package game

import "strings"

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
		result.WriteString(filename)
		result.WriteString("\n")
	}
	return result.String()
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
