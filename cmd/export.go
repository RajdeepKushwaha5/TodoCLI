package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"todo-cli/internal/todo"
)

var (
	exportFormat string
	exportFile   string
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export tasks to file",
	Long: `Export your tasks to various file formats.

Supported formats:
  - csv: Comma-separated values
  - txt: Plain text format

Examples:
  todo export --format=csv --file=tasks.csv
  todo export --format=txt --file=tasks.txt
  todo export --format=csv                    # Exports to tasks.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set default filename based on format
		if exportFile == "" {
			switch exportFormat {
			case "csv":
				exportFile = "tasks.csv"
			case "txt":
				exportFile = "tasks.txt"
			default:
				exportFile = "tasks.txt"
			}
		}

		// Get all tasks
		allTasks := manager.ListTasks(todo.FilterOptions{
			ShowCompleted: true,
			ShowPending:   true,
		})

		switch exportFormat {
		case "csv":
			return exportCSV(allTasks, exportFile)
		case "txt":
			return exportTXT(allTasks, exportFile)
		default:
			return fmt.Errorf("unsupported format '%s'. Supported formats: csv, txt", exportFormat)
		}
	},
}

// exportCSV exports tasks to CSV format
func exportCSV(tasks []*todo.Task, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Title", "Completed", "Priority", "Due Date", "Created At", "Updated At"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write tasks
	for _, task := range tasks {
		record := []string{
			strconv.Itoa(task.ID),
			task.Title,
			strconv.FormatBool(task.Completed),
			string(task.Priority),
			"",
			task.CreatedAt.Format("2006-01-02 15:04:05"),
			task.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if task.DueDate != nil {
			record[4] = task.DueDate.Format("2006-01-02 15:04:05")
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	fmt.Printf("📄 Tasks exported to %s (%d tasks)\n", filename, len(tasks))
	return nil
}

// exportTXT exports tasks to plain text format
func exportTXT(tasks []*todo.Task, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "Todo List Export\n")
	fmt.Fprintf(file, "================\n\n")
	fmt.Fprintf(file, "Total tasks: %d\n\n", len(tasks))

	for _, task := range tasks {
		status := "PENDING"
		if task.Completed {
			status = "COMPLETED"
		}

		fmt.Fprintf(file, "[%d] %s\n", task.ID, task.Title)
		fmt.Fprintf(file, "    Status: %s\n", status)
		fmt.Fprintf(file, "    Priority: %s\n", strings.ToUpper(string(task.Priority)))
		
		if task.DueDate != nil {
			fmt.Fprintf(file, "    Due: %s\n", task.DueDate.Format("2006-01-02 15:04"))
		}
		
		fmt.Fprintf(file, "    Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04"))
		
		if task.Completed {
			fmt.Fprintf(file, "    Completed: %s\n", task.UpdatedAt.Format("2006-01-02 15:04"))
		}
		
		fmt.Fprintf(file, "\n")
	}

	fmt.Printf("📄 Tasks exported to %s (%d tasks)\n", filename, len(tasks))
	return nil
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Add flags
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "txt", "Export format (csv, txt)")
	exportCmd.Flags().StringVarP(&exportFile, "file", "o", "", "Output filename")
}