package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Todo structure
type Todo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"createdAt"`
}

// CodeFile structure
type CodeFile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

// CodeFolder structure
type CodeFolder struct {
	ID    string     `json:"id"`
	Name  string     `json:"name"`
	Files []CodeFile `json:"files"`
}

// App struct
type App struct {
	ctx          context.Context
	todos        []Todo
	codeFolders  []CodeFolder
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		todos:       make([]Todo, 0),
		codeFolders: make([]CodeFolder, 0),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.loadTodosFromFile()
	a.loadCodeTemplatesFromFile()
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

// GetTodos returns all todos
func (a *App) GetTodos() []Todo {
	return a.todos
}

// AddTodo adds a new todo
func (a *App) AddTodo(todo Todo) error {
	a.todos = append(a.todos, todo)
	return a.saveTodosToFile()
}

// DeleteTodo deletes a todo by ID
func (a *App) DeleteTodo(id string) error {
	for i, todo := range a.todos {
		if todo.ID == id {
			a.todos = append(a.todos[:i], a.todos[i+1:]...)
			break
		}
	}
	return a.saveTodosToFile()
}

// saveTodosToFile saves todos to file
func (a *App) saveTodosToFile() error {
	configDir := a.GetConfigDir()
	filePath := filepath.Join(configDir, "todos.json")

	data, err := json.Marshal(a.todos)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// loadTodosFromFile loads todos from file
func (a *App) loadTodosFromFile() error {
	configDir := a.GetConfigDir()
	filePath := filepath.Join(configDir, "todos.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		// File doesn't exist, initialize empty
		a.todos = make([]Todo, 0)
		return nil
	}

	return json.Unmarshal(data, &a.todos)
}

// Code Template APIs

// GetCodeFolders returns all code folders
func (a *App) GetCodeFolders() []CodeFolder {
	return a.codeFolders
}

// SaveCodeTemplates saves code templates to file
func (a *App) SaveCodeTemplates(jsonData string) error {
	configDir := a.GetConfigDir()
	filePath := filepath.Join(configDir, "codeTemplates.json")
	return os.WriteFile(filePath, []byte(jsonData), 0644)
}

// LoadCodeTemplates loads code templates from file
func (a *App) LoadCodeTemplates() (string, error) {
	configDir := a.GetConfigDir()
	filePath := filepath.Join(configDir, "codeTemplates.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		// Return empty array if file doesn't exist
		return "[]", nil
	}

	return string(data), nil
}

// ReadFileContent reads file content
func (a *App) ReadFileContent(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// OpenWithDefaultApp opens file with default application
func (a *App) OpenWithDefaultApp(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	default: // linux
		cmd = exec.Command("xdg-open", filePath)
	}

	return cmd.Start()
}

// SelectCodeFile opens file dialog for selecting code files
func (a *App) SelectCodeFile() (string, error) {
	filePath, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "选择代码文件",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "代码文件", Pattern: "*.py;*.r;*.R;*.exe;*.js;*.ts;*.java;*.cpp;*.c;*.go"},
			{DisplayName: "所有文件", Pattern: "*.*"},
		},
	})
	return filePath, err
}

// loadCodeTemplatesFromFile loads code templates from file
func (a *App) loadCodeTemplatesFromFile() error {
	configDir := a.GetConfigDir()
	filePath := filepath.Join(configDir, "codeTemplates.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		// File doesn't exist, initialize empty
		a.codeFolders = make([]CodeFolder, 0)
		return nil
	}

	return json.Unmarshal(data, &a.codeFolders)
}
