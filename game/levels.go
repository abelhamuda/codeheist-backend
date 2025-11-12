package game

func initializeLevels() map[int]*Level {
	return map[int]*Level{
		0: {
			ID:          0,
			Title:       "The Beginning",
			Description: "The password for the next level is stored in a file called readme",
			WelcomeMsg:  "Welcome to CodeHeist! Your first mission is to find the password in the readme file.",
			Filesystem: map[string]interface{}{
				"readme": "bandit1{NH2SXQwcBdpmTEzi3bvBHMM9H66vVXjL}",
			},
			Solution: "bandit1{NH2SXQwcBdpmTEzi3bvBHMM9H66vVXjL}",
			Hint:     "Try using 'ls' to see files, then 'cat readme' to read it",
		},
		1: {
			ID:          1,
			Title:       "The Dash File",
			Description: "The password for the next level is stored in a file called -",
			WelcomeMsg:  "Good job! Now find the password in a file named with just a dash.",
			Filesystem: map[string]interface{}{
				"-":      "bandit2{rRGizSaX8Mk1RTb1CNQoXTcYZWU6lgzi}",
				"readme": "This is a decoy file. The real password is in the file named '-'",
			},
			Solution: "bandit2{rRGizSaX8Mk1RTb1CNQoXTcYZWU6lgzi}",
			Hint:     "Files with special names need special handling. Try 'cat ./-' or 'cat -- -'",
		},
		2: {
			ID:          2,
			Title:       "Spaces in Filename",
			Description: "The password is in a file with spaces in its name",
			WelcomeMsg:  "Now dealing with filenames that contain spaces.",
			Filesystem: map[string]interface{}{
				"file with spaces.txt": "bandit3{6zPeziLdR2RKNdNYFNb6nVCKzphlXHBM}",
				"normal_file.txt":      "This is not the password file",
			},
			Solution: "bandit3{6zPeziLdR2RKNdNYFNb6nVCKzphlXHBM}",
			Hint:     "Use quotes around filenames with spaces: 'cat \"file with spaces.txt\"'",
		},
		3: {
			ID:          3,
			Title:       "Hidden Files",
			Description: "The password is stored in a hidden file",
			WelcomeMsg:  "Some files are hidden from normal view. Can you find them?",
			Filesystem: map[string]interface{}{
				".hidden":     "bandit4{UV1B0aH0fUHJ6dGVyIFN0cm9uZ1Bhc3N3MHJk}",
				"visible.txt": "This file is visible but not useful",
			},
			Solution: "bandit4{UV1B0aH0fUHJ6dGVyIFN0cm9uZ1Bhc3N3MHJk}",
			Hint:     "Hidden files start with a dot. Use 'ls -a' to see all files",
		},
		// --- NEW CHALLENGING LEVELS ---
		4: {
			ID:          4,
			Title:       "File Permissions",
			Description: "The password is in a file you don't have permission to read",
			WelcomeMsg:  "Sometimes files are protected. You need the right permissions to access them.",
			Filesystem: map[string]interface{}{
				"secret.txt":   "bandit5{FilePermissionMaster2024}",
				"readable.txt": "This file is readable but not helpful",
				".permissions": "Try: chmod 700 secret.txt",
			},
			Solution: "bandit5{FilePermissionMaster2024}",
			Hint:     "Some files need permission changes. Use 'chmod' to modify file permissions.",
		},
		5: {
			ID:          5,
			Title:       "Grep Master",
			Description: "Find the password hidden in a large text file",
			WelcomeMsg:  "Now you need to search through content. The password is somewhere in a large file.",
			Filesystem: map[string]interface{}{
				"data.log": `Server started
User login: admin
Error: connection timeout
Password: bandit6{GrepNinja2024}
Debug: memory allocation
Warning: disk space low
Info: backup completed
Error: null pointer exception
User logout: admin
Server stopped`,
				"notes.txt": "The log file contains important information among all the noise.",
			},
			Solution: "bandit6{GrepNinja2024}",
			Hint:     "Use 'grep' to search for patterns in files. Try: grep Password data.log",
		},
		6: {
			ID:          6,
			Title:       "Binary Detective",
			Description: "Extract text from a binary file",
			WelcomeMsg:  "Some files aren't plain text. You'll need special tools to extract readable content.",
			Filesystem: map[string]interface{}{
				"binary.data": "← Binary data → bandit7{BinaryHunter} ← More binary data →",
				"hint.txt":    "Sometimes binary files contain readable strings...",
			},
			Solution: "bandit7{BinaryHunter}",
			Hint:     "Use 'strings' command to extract readable text from binary files.",
		},
		7: {
			ID:          7,
			Title:       "The Maze of Directories",
			Description: "Find the password hidden deep in directory structures",
			WelcomeMsg:  "The filesystem can be complex. Navigate through directories to find what you need.",
			Filesystem: map[string]interface{}{
				"dir1/file1.txt":        "Not here",
				"dir1/dir2/notes.txt":   "Keep looking",
				"dir1/dir2/dir3/secret": "bandit8{DirectoryExplorer}",
				"dir1/decoy.txt":        "Wrong path",
				"dir4/another.txt":      "Dead end",
			},
			Solution: "bandit8{DirectoryExplorer}",
			Hint:     "Use 'find' command to search through directories recursively.",
		},
		8: {
			ID:          8,
			Title:       "Environment Secrets",
			Description: "The password is stored in an environment variable",
			WelcomeMsg:  "Systems often store secrets in environment variables. Can you find them?",
			Filesystem: map[string]interface{}{
				"config.txt": "SECRET_KEY=bandit9{EnvVariableMaster}",
				"script.sh":  "#!/bin/bash\necho $SECRET_KEY",
				"readme.md":  "Check the environment variables...",
			},
			Solution: "bandit9{EnvVariableMaster}",
			Hint:     "Use 'echo' to display environment variables. The password might be in a variable.",
		},
		9: {
			ID:          9,
			Title:       "The Encoded Secret",
			Description: "Decode a base64 encoded password",
			WelcomeMsg:  "Sometimes secrets are encoded to hide them in plain sight. Can you decode it?",
			Filesystem: map[string]interface{}{
				"encoded.txt": "YmFuZGl0MTB7QmFzZTY0RGVjb2RlckF3ZXNvbWV9", // base64 for "bandit10{Base64DecoderAwesome}"
				"hint.txt":    "This looks like base64 encoding...",
			},
			Solution: "bandit10{Base64DecoderAwesome}",
			Hint:     "Use 'base64 -d' to decode base64 encoded text.",
		},
		// Add more levels as needed
	}
}
