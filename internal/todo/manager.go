package todo

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"todo-cli/storage"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidID    = errors.New("invalid task ID")
)

// FilterOptions defines options for filtering tasks
type FilterOptions struct {
	ShowCompleted bool
	ShowPending   bool
	Priority      Priority
	Search        string
	SortBy        string // "id", "priority", "due", "created"
}

// Manager handles all task operations
type Manager struct {
	storage *storage.FileStorage
	tasks   []*Task
	nextID  int
}

// NewManager creates a new task manager
func NewManager(storagePath ...string) *Manager {
	return &Manager{
		storage: storage.NewFileStorage(storagePath...),
		tasks:   []*Task{},
		nextID:  1,
	}
}

// LoadTasks loads tasks from storage
func (m *Manager) LoadTasks() error {
	storageTasks, nextID, err := m.storage.LoadTasks()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}
	
	// Convert storage tasks to domain tasks
	m.tasks = make([]*Task, len(storageTasks))
	for i, st := range storageTasks {
		m.tasks[i] = &Task{
			ID:        st.ID,
			Title:     st.Title,
			Completed: st.Completed,
			DueDate:   st.DueDate,
			Priority:  Priority(st.Priority),
			CreatedAt: st.CreatedAt,
			UpdatedAt: st.UpdatedAt,
		}
	}
	
	m.nextID = nextID
	return nil
}

// SaveTasks saves tasks to storage
func (m *Manager) SaveTasks() error {
	// Convert domain tasks to storage tasks
	storageTasks := make([]*storage.Task, len(m.tasks))
	for i, t := range m.tasks {
		storageTasks[i] = &storage.Task{
			ID:        t.ID,
			Title:     t.Title,
			Completed: t.Completed,
			DueDate:   t.DueDate,
			Priority:  storage.Priority(t.Priority),
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		}
	}
	
	if err := m.storage.SaveTasks(storageTasks, m.nextID); err != nil {
		return fmt.Errorf("failed to save tasks: %w", err)
	}
	return nil
}

// AddTask adds a new task
func (m *Manager) AddTask(title string, priority Priority, dueDate *time.Time) (*Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("task title cannot be empty")
	}

	task := NewTask(m.nextID, strings.TrimSpace(title))
	
	if priority != "" {
		if !ValidatePriority(string(priority)) {
			return nil, fmt.Errorf("invalid priority: %s (valid options: low, medium, high)", priority)
		}
		task.SetPriority(priority)
	}
	
	if dueDate != nil {
		task.SetDueDate(*dueDate)
	}

	m.tasks = append(m.tasks, task)
	m.nextID++

	if err := m.SaveTasks(); err != nil {
		return nil, err
	}

	return task, nil
}

// GetTask retrieves a task by ID
func (m *Manager) GetTask(id int) (*Task, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	for _, task := range m.tasks {
		if task.ID == id {
			return task, nil
		}
	}
	
	return nil, ErrTaskNotFound
}

// CompleteTask marks a task as completed
func (m *Manager) CompleteTask(id int) (*Task, error) {
	task, err := m.GetTask(id)
	if err != nil {
		return nil, err
	}

	if task.Completed {
		return nil, fmt.Errorf("task %d is already completed", id)
	}

	task.Complete()
	
	if err := m.SaveTasks(); err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask removes a task
func (m *Manager) DeleteTask(id int) (*Task, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	for i, task := range m.tasks {
		if task.ID == id {
			// Remove task from slice
			deletedTask := *task
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			
			if err := m.SaveTasks(); err != nil {
				return nil, err
			}
			
			return &deletedTask, nil
		}
	}
	
	return nil, ErrTaskNotFound
}

// ListTasks returns filtered and sorted tasks
func (m *Manager) ListTasks(filter FilterOptions) []*Task {
	var filteredTasks []*Task

	for _, task := range m.tasks {
		// Apply completion filter
		if filter.ShowCompleted && !filter.ShowPending {
			if !task.Completed {
				continue
			}
		} else if filter.ShowPending && !filter.ShowCompleted {
			if task.Completed {
				continue
			}
		}

		// Apply priority filter
		if filter.Priority != "" {
			if task.Priority != filter.Priority {
				continue
			}
		}

		// Apply search filter
		if filter.Search != "" {
			searchLower := strings.ToLower(filter.Search)
			titleLower := strings.ToLower(task.Title)
			if !strings.Contains(titleLower, searchLower) {
				continue
			}
		}

		filteredTasks = append(filteredTasks, task)
	}

	// Sort tasks
	m.sortTasks(filteredTasks, filter.SortBy)

	return filteredTasks
}

// sortTasks sorts the task slice based on the specified criteria
func (m *Manager) sortTasks(tasks []*Task, sortBy string) {
	switch sortBy {
	case "priority":
		sort.Slice(tasks, func(i, j int) bool {
			priorityOrder := map[Priority]int{
				PriorityHigh:   3,
				PriorityMedium: 2,
				PriorityLow:    1,
			}
			return priorityOrder[tasks[i].Priority] > priorityOrder[tasks[j].Priority]
		})
	case "due":
		sort.Slice(tasks, func(i, j int) bool {
			// Tasks without due dates go to the end
			if tasks[i].DueDate == nil && tasks[j].DueDate == nil {
				return tasks[i].ID < tasks[j].ID
			}
			if tasks[i].DueDate == nil {
				return false
			}
			if tasks[j].DueDate == nil {
				return true
			}
			return tasks[i].DueDate.Before(*tasks[j].DueDate)
		})
	case "created":
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		})
	default: // "id" or any other value
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].ID < tasks[j].ID
		})
	}
}

// GetStats returns statistics about tasks
func (m *Manager) GetStats() map[string]int {
	stats := map[string]int{
		"total":     len(m.tasks),
		"completed": 0,
		"pending":   0,
		"overdue":   0,
		"high":      0,
		"medium":    0,
		"low":       0,
	}

	for _, task := range m.tasks {
		if task.Completed {
			stats["completed"]++
		} else {
			stats["pending"]++
			if task.IsOverdue() {
				stats["overdue"]++
			}
		}

		switch task.Priority {
		case PriorityHigh:
			stats["high"]++
		case PriorityMedium:
			stats["medium"]++
		case PriorityLow:
			stats["low"]++
		}
	}

	return stats
}

// GetStoragePath returns the current storage file path
func (m *Manager) GetStoragePath() string {
	return m.storage.GetFilePath()
}

// BackupTasks creates a backup of tasks
func (m *Manager) BackupTasks() error {
	return m.storage.BackupTasks()
}