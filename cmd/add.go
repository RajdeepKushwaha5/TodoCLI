package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"todo-cli/internal/todo"
)

var (
	addPriority string
	addDueDate  string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [task title]",
	Short: "Add a new task",
	Long: `Add a new task to your todo list.

You can specify priority and due date for better organization.

Examples:
  todo add "Buy groceries"
  todo add "Finish project" --priority=high
  todo add "Meeting with team" --due=2025-10-05
  todo add "Complete assignment" --priority=medium --due=2025-10-10`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Join all arguments to form the task title
		title := ""
		for i, arg := range args {
			if i > 0 {
				title += " "
			}
			title += arg
		}

		// Parse priority
		var priority todo.Priority
		if addPriority != "" {
			if !todo.ValidatePriority(addPriority) {
				return fmt.Errorf("invalid priority '%s'. Valid options: low, medium, high", addPriority)
			}
			priority = todo.Priority(addPriority)
		}

		// Parse due date
		var dueDate *time.Time
		if addDueDate != "" {
			parsedDate, err := time.Parse("2006-01-02", addDueDate)
			if err != nil {
				// Try parsing with time
				parsedDate, err = time.Parse("2006-01-02 15:04", addDueDate)
				if err != nil {
					return fmt.Errorf("invalid due date format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM")
				}
			}
			dueDate = &parsedDate
		}

		// Add the task
		task, err := manager.AddTask(title, priority, dueDate)
		if err != nil {
			return fmt.Errorf("failed to add task: %w", err)
		}

		// Display success message
		fmt.Printf("✅ Task added successfully!\n")
		fmt.Printf("   ID: %d\n", task.ID)
		fmt.Printf("   Title: %s\n", task.Title)
		fmt.Printf("   Priority: %s\n", task.Priority)
		if task.DueDate != nil {
			fmt.Printf("   Due: %s\n", task.DueDate.Format("2006-01-02 15:04"))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Add flags
	addCmd.Flags().StringVarP(&addPriority, "priority", "p", "", "Task priority (low, medium, high)")
	addCmd.Flags().StringVarP(&addDueDate, "due", "d", "", "Due date (YYYY-MM-DD or YYYY-MM-DD HH:MM)")
}