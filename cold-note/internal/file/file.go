package file

import (
	"fmt"
	"os"
	"strings"
)

func GetDirectories(obsidianDir string) ([]string, error) {
	entries, err := os.ReadDir(obsidianDir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			if len(name) > 0 && (name[0] >= '0' && name[0] <= '9') {
				dirs = append(dirs, name)
			}
		}
	}

	if len(dirs) == 0 {
		return nil, fmt.Errorf("no directories found in %s", obsidianDir)
	}

	return dirs, nil
}

func GetTemplates(templatesDir string) ([]string, error) {
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, fmt.Errorf("error reading templates directory: %w", err)
	}

	var templates []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			templates = append(templates, entry.Name())
		}
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("no template files found in %s", templatesDir)
	}

	return templates, nil
}

func ReadTemplateContent(templatePath string) (string, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error reading template: %w", err)
	}
	return string(content), nil
}
