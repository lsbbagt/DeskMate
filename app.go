package main

import (
	"context"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetConfigDir returns the config directory path
func (a *App) GetConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	configDir := filepath.Join(homeDir, ".statbox-v2")
	
	// Create directory if not exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
	}
	
	return configDir
}

// SaveBookmarks saves bookmarks to file
func (a *App) SaveBookmarks(jsonData string) error {
	configDir := a.GetConfigDir()
	filePath := filepath.Join(configDir, "bookmarks.json")
	return os.WriteFile(filePath, []byte(jsonData), 0644)
}

// LoadBookmarks loads bookmarks from file
func (a *App) LoadBookmarks() (string, error) {
	configDir := a.GetConfigDir()
	filePath := filepath.Join(configDir, "bookmarks.json")
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		// Return empty array if file doesn't exist
		return "[]", nil
	}
	
	return string(data), nil
}
