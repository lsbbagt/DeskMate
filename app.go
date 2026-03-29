package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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

// MarkdownFile structure
type MarkdownFile struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

// MarkdownFolder structure
type MarkdownFolder struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Files []MarkdownFile `json:"files"`
}

// App struct
type App struct {
	ctx            context.Context
	todos          []Todo
	codeFolders    []CodeFolder
	markdownFolders []MarkdownFolder
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		todos:           make([]Todo, 0),
		codeFolders:     make([]CodeFolder, 0),
		markdownFolders: make([]MarkdownFolder, 0),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.loadTodosFromFile()
	a.loadCodeTemplatesFromFile()
	a.loadMarkdownFilesFromFile()
	
	// Sync with file system
	a.SyncCodeTemplates()
	a.SyncMarkdownFiles()
}

// GetConfigDir returns the config directory path
func (a *App) GetConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	configDir := filepath.Join(homeDir, ".deskmate")

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

// UpdateTodo updates an existing todo
func (a *App) UpdateTodo(updatedTodo Todo) error {
	for i, todo := range a.todos {
		if todo.ID == updatedTodo.ID {
			a.todos[i] = updatedTodo
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

// saveCodeTemplatesToFile saves code templates to file (internal)
func (a *App) saveCodeTemplatesToFile() error {
	configDir := a.GetConfigDir()
	filePath := filepath.Join(configDir, "codeTemplates.json")

	data, err := json.Marshal(a.codeFolders)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
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
			{DisplayName: "代码文件", Pattern: "*.py;*.r;*.R;*.jl;*.julia;*.c;*.cpp;*.go;*.html;*.m;*.txt"},
			{DisplayName: "所有文件", Pattern: "*.*"},
		},
	})
	return filePath, err
}

// ImportCodeFile copies a code file to the data folder
func (a *App) ImportCodeFile(sourcePath string, folderName string) (string, error) {
	// Get file name
	fileName := filepath.Base(sourcePath)

	// Create folder directory with folder name
	configDir := a.GetConfigDir()
	targetDir := filepath.Join(configDir, "code-templates", folderName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}

	// Copy file
	targetPath := filepath.Join(targetDir, fileName)
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return "", err
	}

	return targetPath, nil
}

// DeleteCodeFile deletes a code file from the data folder
func (a *App) DeleteCodeFile(filePath string) error {
	return os.Remove(filePath)
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

// SyncCodeTemplates syncs code templates with actual file system
func (a *App) SyncCodeTemplates() error {
	configDir := a.GetConfigDir()
	codeTemplatesDir := filepath.Join(configDir, "code-templates")

	// Create directory if not exists
	if _, err := os.Stat(codeTemplatesDir); os.IsNotExist(err) {
		os.MkdirAll(codeTemplatesDir, 0755)
		a.codeFolders = make([]CodeFolder, 0)
		return nil
	}

	// Read actual folders from file system
	actualFolders := make(map[string][]CodeFile)
	
	entries, err := os.ReadDir(codeTemplatesDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			folderName := entry.Name()
			folderPath := filepath.Join(codeTemplatesDir, folderName)
			
			// Read files in this folder
			fileEntries, err := os.ReadDir(folderPath)
			if err != nil {
				continue
			}

			files := make([]CodeFile, 0)
			for _, fileEntry := range fileEntries {
				if !fileEntry.IsDir() {
					fileName := fileEntry.Name()
					ext := strings.ToLower(filepath.Ext(fileName))
					if len(ext) > 0 {
						ext = ext[1:] // Remove the dot
					}

					files = append(files, CodeFile{
						ID:   filepath.Join(folderName, fileName),
						Name: fileName,
						Path: filepath.Join(folderPath, fileName),
						Type: ext,
					})
				}
			}
			actualFolders[folderName] = files
		}
	}

	// Merge with existing config
	mergedFolders := make([]CodeFolder, 0)
	
	// Create a map of existing folders
	existingFolders := make(map[string]CodeFolder)
	for _, folder := range a.codeFolders {
		existingFolders[folder.Name] = folder
	}

	// Process actual folders
	for folderName, files := range actualFolders {
		if existingFolder, exists := existingFolders[folderName]; exists {
			// Update existing folder's files
			existingFolder.Files = files
			mergedFolders = append(mergedFolders, existingFolder)
		} else {
			// Create new folder
			mergedFolders = append(mergedFolders, CodeFolder{
				ID:    folderName,
				Name:  folderName,
				Files: files,
			})
		}
	}

	a.codeFolders = mergedFolders
	return a.saveCodeTemplatesToFile()
}

// Markdown File APIs

// SaveMarkdownFiles saves markdown files to file
func (a *App) SaveMarkdownFiles(jsonData string) error {
	configDir := a.GetConfigDir()
	filePath := filepath.Join(configDir, "markdownFiles.json")
	return os.WriteFile(filePath, []byte(jsonData), 0644)
}

// saveMarkdownFilesToFile saves markdown files to file (internal)
func (a *App) saveMarkdownFilesToFile() error {
	configDir := a.GetConfigDir()
	filePath := filepath.Join(configDir, "markdownFiles.json")

	data, err := json.Marshal(a.markdownFolders)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// LoadMarkdownFiles loads markdown files from file
func (a *App) LoadMarkdownFiles() (string, error) {
	configDir := a.GetConfigDir()
	filePath := filepath.Join(configDir, "markdownFiles.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		// Return empty array if file doesn't exist
		return "[]", nil
	}

	return string(data), nil
}

// SaveMarkdownContent saves markdown content to a file
func (a *App) SaveMarkdownContent(filePath string, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

// SelectMarkdownFile opens file dialog for selecting markdown files
func (a *App) SelectMarkdownFile() (string, error) {
	filePath, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "选择 Markdown 文件",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Markdown 文件", Pattern: "*.md;*.markdown"},
			{DisplayName: "所有文件", Pattern: "*.*"},
		},
	})
	return filePath, err
}

// ImportMarkdownFile copies a markdown file to the data folder
func (a *App) ImportMarkdownFile(sourcePath string, folderName string) (string, error) {
	// Get file name
	fileName := filepath.Base(sourcePath)

	// Create folder directory with folder name
	configDir := a.GetConfigDir()
	targetDir := filepath.Join(configDir, "markdown-files", folderName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}

	// Copy file
	targetPath := filepath.Join(targetDir, fileName)
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return "", err
	}

	return targetPath, nil
}

// DeleteMarkdownFile deletes a markdown file from the data folder
func (a *App) DeleteMarkdownFile(filePath string) error {
	return os.Remove(filePath)
}

// DeleteCodeFolder deletes a code folder and all its files
func (a *App) DeleteCodeFolder(folderName string) error {
	configDir := a.GetConfigDir()
	folderPath := filepath.Join(configDir, "code-templates", folderName)
	return os.RemoveAll(folderPath)
}

// DeleteMarkdownFolder deletes a markdown folder and all its files
func (a *App) DeleteMarkdownFolder(folderName string) error {
	configDir := a.GetConfigDir()
	folderPath := filepath.Join(configDir, "markdown-files", folderName)
	return os.RemoveAll(folderPath)
}

// CreateMarkdownFile creates a new markdown file in the data folder
func (a *App) CreateMarkdownFile(folderID string, fileName string) (string, error) {
	// Create folder directory
	configDir := a.GetConfigDir()
	targetDir := filepath.Join(configDir, "markdown-files", folderID)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}

	// Create empty file
	targetPath := filepath.Join(targetDir, fileName)
	if err := os.WriteFile(targetPath, []byte(""), 0644); err != nil {
		return "", err
	}

	return targetPath, nil
}

// loadMarkdownFilesFromFile loads markdown files from file
func (a *App) loadMarkdownFilesFromFile() error {
	configDir := a.GetConfigDir()
	filePath := filepath.Join(configDir, "markdownFiles.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		// File doesn't exist, initialize empty
		a.markdownFolders = make([]MarkdownFolder, 0)
		return nil
	}

	return json.Unmarshal(data, &a.markdownFolders)
}

// SyncMarkdownFiles syncs markdown files with actual file system
func (a *App) SyncMarkdownFiles() error {
	configDir := a.GetConfigDir()
	markdownFilesDir := filepath.Join(configDir, "markdown-files")

	// Create directory if not exists
	if _, err := os.Stat(markdownFilesDir); os.IsNotExist(err) {
		os.MkdirAll(markdownFilesDir, 0755)
		a.markdownFolders = make([]MarkdownFolder, 0)
		return nil
	}

	// Read actual folders from file system
	actualFolders := make(map[string][]MarkdownFile)
	
	entries, err := os.ReadDir(markdownFilesDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			folderName := entry.Name()
			folderPath := filepath.Join(markdownFilesDir, folderName)
			
			// Read files in this folder
			fileEntries, err := os.ReadDir(folderPath)
			if err != nil {
				continue
			}

			files := make([]MarkdownFile, 0)
			for _, fileEntry := range fileEntries {
				if !fileEntry.IsDir() {
					fileName := fileEntry.Name()
					
					files = append(files, MarkdownFile{
						ID:   filepath.Join(folderName, fileName),
						Name: fileName,
						Path: filepath.Join(folderPath, fileName),
					})
				}
			}
			actualFolders[folderName] = files
		}
	}

	// Merge with existing config
	mergedFolders := make([]MarkdownFolder, 0)
	
	// Create a map of existing folders
	existingFolders := make(map[string]MarkdownFolder)
	for _, folder := range a.markdownFolders {
		existingFolders[folder.Name] = folder
	}

	// Process actual folders
	for folderName, files := range actualFolders {
		if existingFolder, exists := existingFolders[folderName]; exists {
			// Update existing folder's files
			existingFolder.Files = files
			mergedFolders = append(mergedFolders, existingFolder)
		} else {
			// Create new folder
			mergedFolders = append(mergedFolders, MarkdownFolder{
				ID:    folderName,
				Name:  folderName,
				Files: files,
			})
		}
	}

	a.markdownFolders = mergedFolders
	return a.saveMarkdownFilesToFile()
}

// OpenDataFolder opens the data folder in file explorer
func (a *App) OpenDataFolder() (string, error) {
	configDir := a.GetConfigDir()

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", configDir)
	case "darwin":
		cmd = exec.Command("open", configDir)
	default: // linux
		cmd = exec.Command("xdg-open", configDir)
	}

	err := cmd.Start()
	return configDir, err
}
