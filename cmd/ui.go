package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"todo-cli/internal/todo"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// uiCmd represents the interactive UI command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Launch interactive terminal UI",
	Long: `Launch an interactive terminal-based UI to manage your tasks.

This provides a simple, menu-driven interface to:
  - View all tasks
  - Add new tasks
  - Complete tasks
  - Delete tasks
  - View statistics
  - Export tasks

Examples:
  todo ui              # Launch interactive UI`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInteractiveUI()
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}

// runInteractiveUI starts the interactive terminal UI
func runInteractiveUI() error {
	reader := bufio.NewReader(os.Stdin)

	// Welcome banner
	clearScreen()
	showWelcomeBanner()

	for {
		showMainMenu()

		choice, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			viewTasks()
		case "2":
			addTaskUI(reader)
		case "3":
			completeTaskUI(reader)
		case "4":
			deleteTaskUI(reader)
		case "5":
			viewStatistics()
		case "6":
			exportTasksUI(reader)
		case "7":
			backupTasksUI()
		case "8", "q", "Q":
			showExitMessage()
			return nil
		case "c", "C":
			clearScreen()
		default:
			color.Red("❌ Invalid option. Please try again.")
			pause()
		}
	}
}

// showWelcomeBanner displays the welcome banner
func showWelcomeBanner() {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	fmt.Println()
	cyan.Println("╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║                                                            ║")
	cyan.Println("║              📋  TODO CLI - Task Manager  📋              ║")
	cyan.Println("║                                                            ║")
	cyan.Println("║           Manage Your Tasks Efficiently & Simply          ║")
	cyan.Println("║                                                            ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	yellow.Printf("   📅 Today: %s\n", time.Now().Format("Monday, January 2, 2006"))
	fmt.Println()
}

// showMainMenu displays the main menu
func showMainMenu() {
	green := color.New(color.FgGreen, color.Bold)
	white := color.New(color.FgWhite)

	fmt.Println(strings.Repeat("─", 62))
	green.Println("\n  📋 MAIN MENU")
	fmt.Println(strings.Repeat("─", 62))

	white.Println("\n  1️⃣  View All Tasks")
	white.Println("  2️⃣  Add New Task")
	white.Println("  3️⃣  Complete Task")
	white.Println("  4️⃣  Delete Task")
	white.Println("  5️⃣  View Statistics")
	white.Println("  6️⃣  Export Tasks")
	white.Println("  7️⃣  Backup Tasks")
	white.Println("  8️⃣  Exit (or press 'q')")
	fmt.Println()
	white.Println("  💡 Type 'c' to clear screen")
	fmt.Println(strings.Repeat("─", 62))

	fmt.Print("\n  👉 Choose an option: ")
}

// viewTasks displays all tasks
func viewTasks() {
	clearScreen()

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║                     📋 ALL TASKS                          ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	filter := todo.FilterOptions{
		ShowCompleted: true,
		ShowPending:   true,
		SortBy:        "id",
	}

	tasks := manager.ListTasks(filter)

	if len(tasks) == 0 {
		color.Yellow("\n  ℹ️  No tasks found. Add your first task to get started!\n")
	} else {
		displayTasksUI(tasks)
	}

	pause()
}

// displayTasksUI formats and displays tasks in UI mode
func displayTasksUI(tasks []*todo.Task) {
	for _, task := range tasks {
		// Status indicator
		statusIcon := "⭕"
		statusColor := color.New(color.FgYellow)
		if task.Completed {
			statusIcon = "✅"
			statusColor = color.New(color.FgGreen)
		} else if task.IsOverdue() {
			statusIcon = "🔴"
			statusColor = color.New(color.FgRed)
		}

		// Priority color
		var priorityColor *color.Color
		priorityText := ""
		switch task.Priority {
		case todo.PriorityHigh:
			priorityColor = color.New(color.FgRed, color.Bold)
			priorityText = "HIGH  "
		case todo.PriorityMedium:
			priorityColor = color.New(color.FgYellow)
			priorityText = "MEDIUM"
		case todo.PriorityLow:
			priorityColor = color.New(color.FgGreen)
			priorityText = "LOW   "
		}

		// Display task
		fmt.Print("  ")
		statusColor.Print(statusIcon)
		fmt.Printf(" [%d] ", task.ID)

		if task.Completed {
			color.New(color.Faint).Print(task.Title)
		} else {
			color.New(color.FgWhite, color.Bold).Print(task.Title)
		}

		fmt.Print(" | ")
		priorityColor.Print(priorityText)

		if task.DueDate != nil {
			if task.IsOverdue() && !task.Completed {
				color.New(color.FgRed, color.Bold).Printf(" | ⏰ OVERDUE: %s", task.DueDate.Format("Jan 02"))
			} else {
				color.New(color.FgCyan).Printf(" | 📅 Due: %s", task.DueDate.Format("Jan 02"))
			}
		}

		fmt.Println()

		if task.Completed {
			color.New(color.Faint).Printf("      ✓ Completed on %s\n", task.UpdatedAt.Format("Jan 02, 2006"))
		}

		fmt.Println()
	}
}

// addTaskUI handles adding a new task through UI
func addTaskUI(reader *bufio.Reader) {
	clearScreen()

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║                    ➕ ADD NEW TASK                        ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Get task title
	fmt.Print("  📝 Task Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	if title == "" {
		color.Red("\n  ❌ Task title cannot be empty!")
		pause()
		return
	}

	// Get priority
	fmt.Println("\n  🎯 Priority:")
	fmt.Println("     1. Low")
	fmt.Println("     2. Medium (default)")
	fmt.Println("     3. High")
	fmt.Print("\n  Choose (1-3, or press Enter for default): ")

	priorityChoice, _ := reader.ReadString('\n')
	priorityChoice = strings.TrimSpace(priorityChoice)

	var priority todo.Priority
	switch priorityChoice {
	case "1":
		priority = todo.PriorityLow
	case "3":
		priority = todo.PriorityHigh
	default:
		priority = todo.PriorityMedium
	}

	// Get due date
	fmt.Print("\n  📅 Due Date (YYYY-MM-DD or press Enter to skip): ")
	dueDateStr, _ := reader.ReadString('\n')
	dueDateStr = strings.TrimSpace(dueDateStr)

	var dueDate *time.Time
	if dueDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			color.Red("\n  ⚠️  Invalid date format. Task will be added without due date.")
		} else {
			dueDate = &parsedDate
		}
	}

	// Add the task
	task, err := manager.AddTask(title, priority, dueDate)
	if err != nil {
		color.Red("\n  ❌ Failed to add task: %v", err)
		pause()
		return
	}

	// Success message
	fmt.Println()
	color.Green("  ✅ Task added successfully!")
	fmt.Println()
	fmt.Printf("     ID: %d\n", task.ID)
	fmt.Printf("     Title: %s\n", task.Title)
	fmt.Printf("     Priority: %s\n", strings.ToUpper(string(task.Priority)))
	if task.DueDate != nil {
		fmt.Printf("     Due: %s\n", task.DueDate.Format("2006-01-02"))
	}
	fmt.Println()

	pause()
}

// completeTaskUI handles completing a task through UI
func completeTaskUI(reader *bufio.Reader) {
	clearScreen()

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║                   ✅ COMPLETE TASK                        ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Show pending tasks
	filter := todo.FilterOptions{
		ShowPending:   true,
		ShowCompleted: false,
		SortBy:        "id",
	}

	tasks := manager.ListTasks(filter)

	if len(tasks) == 0 {
		color.Yellow("  ℹ️  No pending tasks to complete!\n")
		pause()
		return
	}

	fmt.Println("  📋 Pending Tasks:\n")
	displayTasksUI(tasks)

	fmt.Print("\n  👉 Enter Task ID to complete (or 0 to cancel): ")

	var taskID int
	_, err := fmt.Fscanf(reader, "%d\n", &taskID)
	if err != nil {
		color.Red("\n  ❌ Invalid input!")
		pause()
		return
	}

	if taskID == 0 {
		color.Yellow("\n  ℹ️  Cancelled.")
		pause()
		return
	}

	task, err := manager.CompleteTask(taskID)
	if err != nil {
		color.Red("\n  ❌ Failed to complete task: %v", err)
		pause()
		return
	}

	fmt.Println()
	color.Green("  ✅ Task completed successfully!")
	fmt.Printf("\n     ID: %d\n", task.ID)
	fmt.Printf("     Title: %s\n", task.Title)
	fmt.Printf("     Completed at: %s\n", task.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()

	pause()
}

// deleteTaskUI handles deleting a task through UI
func deleteTaskUI(reader *bufio.Reader) {
	clearScreen()

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║                    🗑️  DELETE TASK                        ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Show all tasks
	filter := todo.FilterOptions{
		ShowCompleted: true,
		ShowPending:   true,
		SortBy:        "id",
	}

	tasks := manager.ListTasks(filter)

	if len(tasks) == 0 {
		color.Yellow("  ℹ️  No tasks to delete!\n")
		pause()
		return
	}

	fmt.Println("  📋 All Tasks:\n")
	displayTasksUI(tasks)

	fmt.Print("\n  👉 Enter Task ID to delete (or 0 to cancel): ")

	var taskID int
	_, err := fmt.Fscanf(reader, "%d\n", &taskID)
	if err != nil {
		color.Red("\n  ❌ Invalid input!")
		pause()
		return
	}

	if taskID == 0 {
		color.Yellow("\n  ℹ️  Cancelled.")
		pause()
		return
	}

	// Get task details
	task, err := manager.GetTask(taskID)
	if err != nil {
		color.Red("\n  ❌ Failed to find task: %v", err)
		pause()
		return
	}

	// Confirmation
	fmt.Println()
	color.Yellow("  ⚠️  Are you sure you want to delete this task?")
	fmt.Printf("\n     ID: %d\n", task.ID)
	fmt.Printf("     Title: %s\n", task.Title)
	fmt.Print("\n  Type 'yes' to confirm: ")

	confirmation, _ := reader.ReadString('\n')
	confirmation = strings.TrimSpace(strings.ToLower(confirmation))

	if confirmation != "yes" && confirmation != "y" {
		color.Yellow("\n  ℹ️  Deletion cancelled.")
		pause()
		return
	}

	deletedTask, err := manager.DeleteTask(taskID)
	if err != nil {
		color.Red("\n  ❌ Failed to delete task: %v", err)
		pause()
		return
	}

	fmt.Println()
	color.Green("  ✅ Task deleted successfully!")
	fmt.Printf("\n     ID: %d\n", deletedTask.ID)
	fmt.Printf("     Title: %s\n", deletedTask.Title)
	fmt.Println()

	pause()
}

// viewStatistics displays task statistics
func viewStatistics() {
	clearScreen()

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║                   📊 TASK STATISTICS                      ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	stats := manager.GetStats()

	white := color.New(color.FgWhite, color.Bold)

	fmt.Println("  📈 Overview:")
	fmt.Println(strings.Repeat("  ─", 30))
	white.Printf("     Total Tasks:      %d\n", stats["total"])
	color.Green("     ✅ Completed:       %d\n", stats["completed"])
	color.Yellow("     ⭕ Pending:         %d\n", stats["pending"])
	color.Red("     🔴 Overdue:         %d\n", stats["overdue"])
	fmt.Println()

	fmt.Println("  🎯 By Priority:")
	fmt.Println(strings.Repeat("  ─", 30))
	color.New(color.FgRed, color.Bold).Printf("     🔴 High:            %d\n", stats["high"])
	color.New(color.FgYellow).Printf("     🟡 Medium:          %d\n", stats["medium"])
	color.New(color.FgGreen).Printf("     🟢 Low:             %d\n", stats["low"])
	fmt.Println()

	// Completion rate with progress bar
	if stats["total"] > 0 {
		completionRate := float64(stats["completed"]) / float64(stats["total"]) * 100
		fmt.Println("  📊 Progress:")
		fmt.Println(strings.Repeat("  ─", 30))
		white.Printf("     Completion Rate:  %.1f%%\n", completionRate)

		// Progress bar
		barLength := 30
		completed := int(completionRate / 100 * float64(barLength))
		bar := strings.Repeat("█", completed) + strings.Repeat("░", barLength-completed)

		fmt.Print("     ")
		if completionRate >= 75 {
			color.Green("[%s]\n", bar)
		} else if completionRate >= 50 {
			color.Yellow("[%s]\n", bar)
		} else {
			color.Red("[%s]\n", bar)
		}
	}

	fmt.Println()
	color.Cyan("  💾 Storage: %s", manager.GetStoragePath())
	fmt.Println()

	pause()
}

// exportTasksUI handles exporting tasks through UI
func exportTasksUI(reader *bufio.Reader) {
	clearScreen()

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║                   📤 EXPORT TASKS                         ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Println("  📄 Export Format:")
	fmt.Println("     1. CSV (Comma-Separated Values)")
	fmt.Println("     2. TXT (Plain Text)")
	fmt.Print("\n  Choose format (1-2): ")

	formatChoice, _ := reader.ReadString('\n')
	formatChoice = strings.TrimSpace(formatChoice)

	var format, defaultFile string
	switch formatChoice {
	case "1":
		format = "csv"
		defaultFile = "tasks.csv"
	case "2":
		format = "txt"
		defaultFile = "tasks.txt"
	default:
		color.Red("\n  ❌ Invalid format choice!")
		pause()
		return
	}

	fmt.Printf("\n  💾 Filename (press Enter for '%s'): ", defaultFile)
	filename, _ := reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	if filename == "" {
		filename = defaultFile
	}

	// Get all tasks
	allTasks := manager.ListTasks(todo.FilterOptions{
		ShowCompleted: true,
		ShowPending:   true,
	})

	var err error
	switch format {
	case "csv":
		err = exportCSV(allTasks, filename)
	case "txt":
		err = exportTXT(allTasks, filename)
	}

	if err != nil {
		color.Red("\n  ❌ Export failed: %v", err)
	} else {
		fmt.Println()
		color.Green("  ✅ Tasks exported successfully!")
		fmt.Printf("\n     📁 File: %s\n", filename)
		fmt.Printf("     📊 Tasks: %d\n", len(allTasks))
		fmt.Println()
	}

	pause()
}

// backupTasksUI handles backup through UI
func backupTasksUI() {
	clearScreen()

	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("\n╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║                   💾 BACKUP TASKS                         ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	err := manager.BackupTasks()
	if err != nil {
		color.Red("  ❌ Backup failed: %v\n", err)
	} else {
		color.Green("  ✅ Backup created successfully!")
		fmt.Printf("\n     📁 Location: %s.backup\n", manager.GetStoragePath())
		fmt.Println()
	}

	pause()
}

// showExitMessage displays exit message
func showExitMessage() {
	clearScreen()

	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)

	fmt.Println()
	cyan.Println("╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║                                                            ║")
	cyan.Println("║              👋 Thank You for Using Todo CLI!             ║")
	cyan.Println("║                                                            ║")
	cyan.Println("║             Stay Organized & Stay Productive!              ║")
	cyan.Println("║                                                            ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	yellow.Println("  💡 Tip: Run 'todo --help' for command-line options")
	yellow.Println("  🌟 Star us on GitHub: github.com/RajdeepKushwaha5/TodoCLI")
	fmt.Println()
}

// clearScreen clears the terminal screen
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// pause waits for user to press Enter
func pause() {
	fmt.Println()
	color.New(color.FgCyan).Print("  Press Enter to continue...")
	fmt.Scanln()
}
