package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
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

// App struct
type App struct {
	ctx   context.Context
	todos []Todo
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		todos: make([]Todo, 0),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.loadTodosFromFile()
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
