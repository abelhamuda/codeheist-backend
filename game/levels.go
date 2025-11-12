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
		// Add more levels as needed
	}
}
