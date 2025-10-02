package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"todo-cli/internal/todo"
)

var (
	manager *todo.Manager
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A simple and efficient CLI todo manager",
	Long: `Todo CLI is a command-line todo manager that helps you organize your tasks.
You can add, list, complete, and delete tasks with ease.

Examples:
  todo add "Buy groceries" --priority=high --due=2025-10-05
  todo list --completed
  todo complete 1
  todo delete 2`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize the task manager
		manager = todo.NewManager()
		if err := manager.LoadTasks(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load tasks: %v\n", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todo/config.yaml)")
	
	// Add version flag
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetBool("version")
		if version {
			fmt.Println("Todo CLI v1.0.0")
			fmt.Println("A simple and efficient command-line todo manager")
			return
		}
		cmd.Help()
	}
}