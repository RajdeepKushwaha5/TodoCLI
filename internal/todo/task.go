package todo

import (
	"time"
)

// Priority represents task priority levels
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Task represents a single todo item
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	DueDate   *time.Time `json:"due_date,omitempty"`
	Priority  Priority  `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewTask creates a new task with default values
func NewTask(id int, title string) *Task {
	now := time.Now()
	return &Task{
		ID:        id,
		Title:     title,
		Completed: false,
		Priority:  PriorityMedium,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Complete marks a task as completed
func (t *Task) Complete() {
	t.Completed = true
	t.UpdatedAt = time.Now()
}

// SetPriority sets the task priority
func (t *Task) SetPriority(priority Priority) {
	t.Priority = priority
	t.UpdatedAt = time.Now()
}

// SetDueDate sets the task due date
func (t *Task) SetDueDate(dueDate time.Time) {
	t.DueDate = &dueDate
	t.UpdatedAt = time.Now()
}

// IsOverdue checks if the task is overdue
func (t *Task) IsOverdue() bool {
	if t.DueDate == nil || t.Completed {
		return false
	}
	return time.Now().After(*t.DueDate)
}

// ValidatePriority checks if a priority string is valid
func ValidatePriority(priority string) bool {
	switch Priority(priority) {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}