# Todo CLI - A Modern Command-Line Task Manager

[![GitHub](https://img.shields.io/badge/GitHub-Repository-blue?logo=github)](https://github.com/RajdeepKushwaha5/TodoCLI)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A powerful, intuitive, and feature-rich command-line todo application built with Go. Manage your tasks efficiently with priorities, due dates, filtering, and export capabilities.

## 📥 Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/RajdeepKushwaha5/TodoCLI.git
cd TodoCLI

# Build the application
go build -o todo.exe

# Run the application
.\todo.exe --help
```

## ✨ Features

- **Add tasks** with priorities and due dates
- **List tasks** with filtering and sorting options
- **Complete and delete tasks** with confirmation
- **Colored output** for better visual organization
- **Search functionality** to find tasks quickly
- **Export tasks** to CSV or TXT formats
- **Backup functionality** to protect your data
- **Persistent storage** in JSON format
- **Statistics** to track your productivity

## 🚀 Quick Start

### Installing Globally

```bash
# After cloning, install globally
go install

# Use from anywhere (requires GOPATH/bin in PATH)
todo --help
```

## 📖 Usage

### Adding Tasks

```bash
# Add a simple task
todo add "Buy groceries"

# Add a task with priority
todo add "Finish project" --priority=high

# Add a task with due date
todo add "Meeting with team" --due=2025-10-05

# Add a task with both priority and due date
todo add "Complete assignment" --priority=medium --due=2025-10-10
```

### Listing Tasks

```bash
# List all tasks
todo list

# List only completed tasks
todo list --completed

# List only pending tasks
todo list --pending

# Filter by priority
todo list --priority=high

# Search tasks
todo list --search="project"

# Sort tasks
todo list --sort=priority
todo list --sort=due
todo list --sort=created

# Show statistics
todo list --stats
```

### Managing Tasks

```bash
# Complete a task
todo complete 1

# Delete a task (with confirmation)
todo delete 2

# Delete a task without confirmation
todo delete 3 --force
```

### Data Management

```bash
# Export tasks to CSV
todo export --format=csv --file=my-tasks.csv

# Export tasks to TXT
todo export --format=txt --file=my-tasks.txt

# Create a backup
todo backup
```

## 🏗️ Project Structure

```text
todo-cli/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command and app entry
│   ├── add.go             # Add task command
│   ├── list.go            # List tasks command
│   ├── complete.go        # Complete task command
│   ├── delete.go          # Delete task command
│   ├── export.go          # Export tasks command
│   └── backup.go          # Backup command
├── internal/              # Internal packages
│   └── todo/              # Core todo logic
│       ├── task.go        # Task struct and methods
│       └── manager.go     # Task management logic
├── storage/               # Storage layer
│   └── file.go           # JSON file storage
├── main.go               # Application entry point
├── go.mod                # Go module file
└── go.sum                # Go dependencies
```

## 🎨 Output Examples

### Task List

```text
📋 Todo List (3 tasks)
────────────────────────────────────────────────────────────
⭕ [1]    Buy groceries                            HIGH
       Created: Oct 02, 2025

✅ [2]    ✓ Finish Go project                      MEDIUM (Due: 2025-10-05)
       Created: Oct 02, 2025 | Completed: Oct 02, 2025

⭕ [3]    Study Golang                             LOW
       Created: Oct 02, 2025
```

### Statistics

```text
📊 Task Statistics
──────────────────────────────
Total tasks:      3
Completed:        1
Pending:          2
Overdue:          0

By Priority:
  High:           1
  Medium:         1
  Low:            1

Completion rate:  33.3%

Storage location: C:\Users\username\.todo\tasks.json
```

## 🗂️ Data Storage

Tasks are stored in JSON format in:
- **Windows**: `C:\Users\[username]\.todo\tasks.json`
- **Linux/Mac**: `$HOME/.todo/tasks.json`

### Sample JSON Structure
```json
{
  "tasks": [
    {
      "id": 1,
      "title": "Buy groceries",
      "completed": false,
      "due_date": "2025-10-05T00:00:00Z",
      "priority": "high",
      "created_at": "2025-10-02T21:26:51Z",
      "updated_at": "2025-10-02T21:26:51Z"
    }
  ],
  "next_id": 2,
  "last_update": "username"
}
```

## ⚙️ Configuration

### Priority Levels
- `high` - Red color, highest importance
- `medium` - Yellow color, default priority
- `low` - Green color, lowest importance

### Date Formats
- `YYYY-MM-DD` (e.g., 2025-10-05)
- `YYYY-MM-DD HH:MM` (e.g., 2025-10-05 14:30)

### Sort Options
- `id` - Sort by task ID (default)
- `priority` - Sort by priority (high to low)
- `due` - Sort by due date (earliest first)
- `created` - Sort by creation date

## 🔧 Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Color](https://github.com/fatih/color) - Colored terminal output

## 🐛 Troubleshooting

### Common Issues

1. **"todo: command not found"**
   - Ensure GOPATH/bin is in your system PATH
   - Or use the full path to the executable

2. **"Failed to load tasks"**
   - Check file permissions in the .todo directory
   - Ensure the JSON file is not corrupted

3. **"Invalid date format"**
   - Use YYYY-MM-DD or YYYY-MM-DD HH:MM format
   - Example: 2025-10-05 or 2025-10-05 14:30

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository: [TodoCLI](https://github.com/RajdeepKushwaha5/TodoCLI)
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Commit your changes: `git commit -m 'Add some amazing feature'`
5. Push to the branch: `git push origin feature/amazing-feature`
6. Submit a pull request

## �‍💻 Author

**Rajdeep Singh**

- 🐙 GitHub: [@RajdeepKushwaha5](https://github.com/RajdeepKushwaha5)
- 🐦 Twitter/X: [@rajdeeptwts](https://x.com/rajdeeptwts)
- 💼 LinkedIn: [Rajdeep Singh](https://www.linkedin.com/in/rajdeep-singh-b658a833a/)
- 📝 Medium: [@rajdeep01](https://medium.com/@rajdeep01)

## �📝 License

This project is open source and available under the MIT License.

## 🎯 Future Enhancements

- [ ] Task categories/tags
- [ ] Recurring tasks
- [ ] Task dependencies
- [ ] Multiple storage backends
- [ ] Web interface
- [ ] Task time tracking
- [ ] Integration with calendar apps
- [ ] Cloud synchronization
- [ ] Mobile companion app

## ⭐ Show Your Support

If you find this project helpful, please consider:

- ⭐ Starring the repository on [GitHub](https://github.com/RajdeepKushwaha5/TodoCLI)
- 🐛 Reporting any issues you encounter
- 💡 Suggesting new features
- 🤝 Contributing to the codebase

## 📬 Contact

Feel free to reach out if you have any questions or suggestions:

- 📧 Email: rajdeep01@[domain]
- 🐦 Twitter: [@rajdeeptwts](https://x.com/rajdeeptwts)
- 💼 LinkedIn: [Rajdeep Singh](https://www.linkedin.com/in/rajdeep-singh-b658a833a/)
