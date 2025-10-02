package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"todo-cli/internal/todo"
)

var (
	listCompleted bool
	listPending   bool
	listPriority  string
	listSearch    string
	listSort      string
	listStats     bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long: `List your tasks with various filtering and sorting options.

Examples:
  todo list                           # List all tasks
  todo list --completed               # List only completed tasks
  todo list --pending                 # List only pending tasks
  todo list --priority=high           # List only high priority tasks
  todo list --search="project"        # Search for tasks containing "project"
  todo list --sort=priority           # Sort by priority
  todo list --sort=due                # Sort by due date
  todo list --stats                   # Show task statistics`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show statistics if requested
		if listStats {
			return showStats()
		}

		// Build filter options
		filter := todo.FilterOptions{
			ShowCompleted: listCompleted,
			ShowPending:   listPending,
			Search:        listSearch,
			SortBy:        listSort,
		}

		// Validate and set priority filter
		if listPriority != "" {
			if !todo.ValidatePriority(listPriority) {
				return fmt.Errorf("invalid priority '%s'. Valid options: low, medium, high", listPriority)
			}
			filter.Priority = todo.Priority(listPriority)
		}

		// If neither completed nor pending is specified, show all
		if !listCompleted && !listPending {
			filter.ShowCompleted = true
			filter.ShowPending = true
		}

		// Get filtered tasks
		tasks := manager.ListTasks(filter)

		if len(tasks) == 0 {
			fmt.Println("📋 No tasks found matching your criteria.")
			return nil
		}

		// Display tasks
		return displayTasks(tasks)
	},
}

// displayTasks formats and displays the task list
func displayTasks(tasks []*todo.Task) error {
	fmt.Printf("\n📋 Todo List (%d tasks)\n", len(tasks))
	fmt.Println(strings.Repeat("─", 60))

	for _, task := range tasks {
		displayTask(task)
	}

	fmt.Println()
	return nil
}

// displayTask formats and displays a single task
func displayTask(task *todo.Task) {
	// Status icon
	statusIcon := "⭕"
	if task.Completed {
		statusIcon = "✅"
	} else if task.IsOverdue() {
		statusIcon = "🔴"
	}

	// Priority color
	var priorityColor *color.Color
	switch task.Priority {
	case todo.PriorityHigh:
		priorityColor = color.New(color.FgRed, color.Bold)
	case todo.PriorityMedium:
		priorityColor = color.New(color.FgYellow)
	case todo.PriorityLow:
		priorityColor = color.New(color.FgGreen)
	default:
		priorityColor = color.New(color.Reset)
	}

	// Task title (cross out if completed)
	titleColor := color.New(color.Reset)
	title := task.Title
	if task.Completed {
		titleColor = color.New(color.Faint)
		title = "✓ " + title
	}

	// Format ID with padding
	idStr := fmt.Sprintf("[%d]", task.ID)

	// Build the main line
	fmt.Printf("%s %-6s ", statusIcon, idStr)
	titleColor.Printf("%-40s", title)
	priorityColor.Printf(" %s", strings.ToUpper(string(task.Priority)))

	// Due date info
	if task.DueDate != nil {
		dueStr := task.DueDate.Format("2006-01-02")
		if task.IsOverdue() && !task.Completed {
			color.New(color.FgRed, color.Bold).Printf(" (DUE: %s)", dueStr)
		} else {
			color.New(color.FgCyan).Printf(" (Due: %s)", dueStr)
		}
	}

	fmt.Println()

	// Additional info line (created date, etc.)
	createdStr := task.CreatedAt.Format("Jan 02, 2006")
	color.New(color.Faint).Printf("       Created: %s", createdStr)
	
	if task.Completed {
		updatedStr := task.UpdatedAt.Format("Jan 02, 2006")
		color.New(color.Faint).Printf(" | Completed: %s", updatedStr)
	}
	
	fmt.Println()
	fmt.Println()
}

// showStats displays task statistics
func showStats() error {
	stats := manager.GetStats()
	
	fmt.Println("\n📊 Task Statistics")
	fmt.Println(strings.Repeat("─", 30))
	
	fmt.Printf("Total tasks:      %d\n", stats["total"])
	fmt.Printf("Completed:        %d\n", stats["completed"])
	fmt.Printf("Pending:          %d\n", stats["pending"])
	fmt.Printf("Overdue:          %d\n", stats["overdue"])
	fmt.Println()
	
	fmt.Println("By Priority:")
	color.New(color.FgRed).Printf("  High:           %d\n", stats["high"])
	color.New(color.FgYellow).Printf("  Medium:         %d\n", stats["medium"])
	color.New(color.FgGreen).Printf("  Low:            %d\n", stats["low"])
	
	// Completion rate
	if stats["total"] > 0 {
		completionRate := float64(stats["completed"]) / float64(stats["total"]) * 100
		fmt.Printf("\nCompletion rate:  %.1f%%\n", completionRate)
	}
	
	// Storage info
	fmt.Printf("\nStorage location: %s\n", manager.GetStoragePath())
	
	fmt.Println()
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().BoolVarP(&listCompleted, "completed", "c", false, "Show only completed tasks")
	listCmd.Flags().BoolVarP(&listPending, "pending", "p", false, "Show only pending tasks")
	listCmd.Flags().StringVar(&listPriority, "priority", "", "Filter by priority (low, medium, high)")
	listCmd.Flags().StringVarP(&listSearch, "search", "s", "", "Search tasks by title")
	listCmd.Flags().StringVar(&listSort, "sort", "id", "Sort by: id, priority, due, created")
	listCmd.Flags().BoolVar(&listStats, "stats", false, "Show task statistics")
}