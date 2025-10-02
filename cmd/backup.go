package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of your tasks",
	Long: `Create a backup of your tasks file.

This creates a backup file with the .backup extension in the same directory
as your tasks file.

Example:
  todo backup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := manager.BackupTasks(); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}

		fmt.Printf("💾 Backup created successfully!\n")
		fmt.Printf("   Location: %s.backup\n", manager.GetStoragePath())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}