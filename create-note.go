package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func getFilename() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter a filename: ")
		filename, _ := reader.ReadString('\n')
		trimmed := strings.TrimSpace(filename)

		if len(trimmed) > 0 {
			return trimmed
		}

		fmt.Println("Error: The filename cannot be empty. Please insert a proper value.")
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func readTemplate(directory string) (string, error) {

	templatePath := filepath.Join(directory, "template.md") 
	if !fileExists(templatePath) {
		return "", fmt.Errorf("please create template .md file at %s", templatePath)
	}
	content, err := os.ReadFile(templatePath)

	if err != nil {
		return "", fmt.Errorf("error reading template: %w", err)
	}

	return string(content), nil

}

func openFile(directory string, filename string) error {
	if len(filename) == 0 {
		return fmt.Errorf("please insert a filename")
	}
	// Create full path
	fullPath := filepath.Join(directory, filename+".md")

	if fileExists(fullPath) {
		return fmt.Errorf("file already exists: %s", fullPath)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	template, err := readTemplate(directory)
	if err != nil {
		return fmt.Errorf("error getting template %w", err)
	}

	// Write template content
	timestamp := time.Now().Format("200601021504")
	// template := fmt.Sprintf("# \n\n\n\nLinks:\n\n%s\n", timestamp)
	content := strings.ReplaceAll(template, "{{timestamp}}", timestamp)
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	// Open in Neovim with ZenMode
	cmd := exec.Command("nvim", "+ normal ggzzi", fullPath, "-c", ":ZenMode")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func selectDir(zetDir string) (string, error) {
	// Read all entries in the ZETTELKASTEN directory
	entries, err := os.ReadDir(zetDir)
	if err != nil {
		return "", fmt.Errorf("error reading directory: %w", err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			// Check if directory name starts with a number
			if len(name) > 0 && (name[0] >= '0' && name[0] <= '9') {
				dirs = append(dirs, name)
			}
		}
	}
	if len(dirs) == 0 {
		return "", fmt.Errorf("no directories found in %s", zetDir)
	}

	// Print directories with numbers
	fmt.Println("\nAvailable directories:")
	for i, dir := range dirs {
		fmt.Printf("%d: %s\n", i+1, dir)
	}

	// Get user selection
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nSelect directory number: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("error reading input: %w", err)
		}

		// Convert input to number
		var selection int
		_, err = fmt.Sscanf(strings.TrimSpace(input), "%d", &selection)
		if err != nil || selection < 1 || selection > len(dirs) {
			fmt.Printf("Please enter a number between 1 and %d\n", len(dirs))
			continue
		}

		// Return the full path of selected directory
		return filepath.Join(zetDir, dirs[selection-1]), nil
	}
}

func main() {
	var filename string

	// Get zettelkasten directory from environment
	zetDir := os.Getenv("ZETTELKASTEN")
	if zetDir == "" {
		fmt.Println("ZETTELKASTEN environment variable not set")
		os.Exit(1)
	}

	// Process arguments
	args := os.Args[1:]
	switch len(args) {
	case 0:
		filename = getFilename()
	case 1:
		filename = args[0]
	default:
		fmt.Println("Please provide only one filename separated by dashes, without .md extension.")
		fmt.Println("Example: zet my-new-note")
		os.Exit(1)
	}

	// Open file in the inbox directory
	// inboxDir := filepath.Join(zetDir, "0-inbox")

	selectedDir, err := selectDir(zetDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if err := openFile(selectedDir, filename); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
