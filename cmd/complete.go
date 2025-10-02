package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete [task_id]",
	Short: "Mark a task as completed",
	Long: `Mark a task as completed by providing its ID.

Examples:
  todo complete 1      # Mark task with ID 1 as completed
  todo complete 5      # Mark task with ID 5 as completed

You can find task IDs by running: todo list`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse task ID
		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID '%s'. Please provide a valid number", args[0])
		}

		// Complete the task
		task, err := manager.CompleteTask(taskID)
		if err != nil {
			return fmt.Errorf("failed to complete task: %w", err)
		}

		// Display success message
		fmt.Printf("✅ Task completed successfully!\n")
		fmt.Printf("   ID: %d\n", task.ID)
		fmt.Printf("   Title: %s\n", task.Title)
		fmt.Printf("   Completed at: %s\n", task.UpdatedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}