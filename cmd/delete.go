package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	deleteForce bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [task_id]",
	Short: "Delete a task",
	Long: `Delete a task permanently by providing its ID.

Examples:
  todo delete 1        # Delete task with ID 1 (with confirmation)
  todo delete 5 --force # Delete task with ID 5 without confirmation

You can find task IDs by running: todo list

Warning: This action cannot be undone!`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse task ID
		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID '%s'. Please provide a valid number", args[0])
		}

		// Get the task first to show what will be deleted
		task, err := manager.GetTask(taskID)
		if err != nil {
			return fmt.Errorf("failed to find task: %w", err)
		}

		// Confirm deletion unless --force is used
		if !deleteForce {
			fmt.Printf("⚠️  Are you sure you want to delete this task?\n")
			fmt.Printf("   ID: %d\n", task.ID)
			fmt.Printf("   Title: %s\n", task.Title)
			fmt.Printf("   Status: ")
			if task.Completed {
				fmt.Printf("Completed\n")
			} else {
				fmt.Printf("Pending\n")
			}
			fmt.Printf("\nType 'yes' to confirm deletion: ")
			
			var confirmation string
			fmt.Scanln(&confirmation)
			
			if confirmation != "yes" && confirmation != "y" && confirmation != "YES" && confirmation != "Y" {
				fmt.Println("❌ Deletion cancelled.")
				return nil
			}
		}

		// Delete the task
		deletedTask, err := manager.DeleteTask(taskID)
		if err != nil {
			return fmt.Errorf("failed to delete task: %w", err)
		}

		// Display success message
		fmt.Printf("🗑️  Task deleted successfully!\n")
		fmt.Printf("   ID: %d\n", deletedTask.ID)
		fmt.Printf("   Title: %s\n", deletedTask.Title)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	
	// Add flags
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Delete without confirmation")
}