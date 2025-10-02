package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const defaultFileName = "tasks.json"

// Priority represents task priority levels
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Task represents a single todo item for storage
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	DueDate   *time.Time `json:"due_date,omitempty"`
	Priority  Priority  `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FileStorage handles saving and loading tasks to/from JSON files
type FileStorage struct {
	filePath string
}

// TaskList represents the structure stored in JSON
type TaskList struct {
	Tasks      []*Task `json:"tasks"`
	NextID     int     `json:"next_id"`
	LastUpdate string  `json:"last_update"`
}

// NewFileStorage creates a new file storage instance
func NewFileStorage(customPath ...string) *FileStorage {
	var filePath string
	
	if len(customPath) > 0 && customPath[0] != "" {
		filePath = customPath[0]
	} else {
		// Default to user's home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fallback to current directory
			filePath = defaultFileName
		} else {
			filePath = filepath.Join(homeDir, ".todo", defaultFileName)
		}
	}
	
	return &FileStorage{
		filePath: filePath,
	}
}

// ensureDir creates the directory for the storage file if it doesn't exist
func (fs *FileStorage) ensureDir() error {
	dir := filepath.Dir(fs.filePath)
	return os.MkdirAll(dir, 0755)
}

// LoadTasks loads tasks from the JSON file
func (fs *FileStorage) LoadTasks() ([]*Task, int, error) {
	// Ensure directory exists
	if err := fs.ensureDir(); err != nil {
		return nil, 1, fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(fs.filePath); os.IsNotExist(err) {
		// Return empty list if file doesn't exist
		return []*Task{}, 1, nil
	}

	// Read file
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return nil, 1, fmt.Errorf("failed to read file: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		return []*Task{}, 1, nil
	}

	// Parse JSON
	var taskList TaskList
	if err := json.Unmarshal(data, &taskList); err != nil {
		return nil, 1, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Ensure we have a valid NextID
	if taskList.NextID == 0 {
		taskList.NextID = 1
		// Find the highest ID and set NextID accordingly
		for _, task := range taskList.Tasks {
			if task.ID >= taskList.NextID {
				taskList.NextID = task.ID + 1
			}
		}
	}

	return taskList.Tasks, taskList.NextID, nil
}

// SaveTasks saves tasks to the JSON file
func (fs *FileStorage) SaveTasks(tasks []*Task, nextID int) error {
	// Ensure directory exists
	if err := fs.ensureDir(); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create task list structure
	taskList := TaskList{
		Tasks:      tasks,
		NextID:     nextID,
		LastUpdate: fmt.Sprintf("%v", os.Getenv("USER")),
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(taskList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(fs.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetFilePath returns the current file path being used for storage
func (fs *FileStorage) GetFilePath() string {
	return fs.filePath
}

// FileExists checks if the storage file exists
func (fs *FileStorage) FileExists() bool {
	_, err := os.Stat(fs.filePath)
	return !os.IsNotExist(err)
}

// BackupTasks creates a backup of the current tasks file
func (fs *FileStorage) BackupTasks() error {
	if !fs.FileExists() {
		return fmt.Errorf("no tasks file to backup")
	}

	backupPath := fs.filePath + ".backup"
	
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return fmt.Errorf("failed to read tasks file: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	return nil
}