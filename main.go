package main

import (
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Styling definitions
	taskStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#1F6FEB"))
	completedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Strikethrough(true)
	pointerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#1F6FEB"))
	inputStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#1F6FEB"))
	mainStyle      = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

type model struct {
	tasks    []task
	selected int
	creating bool
	input    string
}

type task struct {
	Label     string `json:"label"`
	Completed bool   `json:"completed"`
}

func main() {
	initialTasks := []task{
		{Label: "Wake up sleeepy head 🛏️ 💤", Completed: false},
		{Label: "Brush your teeth 🦷", Completed: false},
		{Label: "Get some ☕", Completed: false},
	}

	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Could not load tasks:", err)
		tasks = initialTasks
	}

	p := tea.NewProgram(model{tasks: tasks})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// TODO: Move into Service file
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.creating {
				m.creating = false
				m.input = ""
				return m, nil
			}
			// Persists tasks to json file
			if err := saveTasks(m.tasks); err != nil {
				return m, tea.Quit
			}
			return m, tea.Quit
		case "down":
			if !m.creating {
				m.selected = (m.selected + 1) % len(m.tasks)
			}
		case "up":
			if !m.creating && m.selected > 0 {
				m.selected--
			} else {
				m.selected = len(m.tasks) - 1
			}
		case "enter":
			if !m.creating {
				m.tasks[m.selected].Completed = !m.tasks[m.selected].Completed
			} else if m.input != "" {
				m.tasks = append(m.tasks, task{Label: m.input, Completed: false})
				m.input = ""
				m.creating = false
			}
		case "n":
			if !m.creating {
				m.creating = true
				m.input = ""
			}
		case "backspace":
			if m.creating && len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			if m.creating {
				m.input += msg.String()
			}
		}
	}

	return m, nil
}

// TODO: Move into View file
func (m model) View() string {
	var s string
	if m.creating {
		s += inputStyle.Render(fmt.Sprintf("Add Task: %s", m.input))
	} else {
		s += "🧙 Mornin', adventurer! ☁️⚡\n\n"
		for i, task := range m.tasks {
			cursor := " "
			if i == m.selected {
				cursor = pointerStyle.Render(">")
			}
			taskText := taskStyle.Render(task.Label)
			if task.Completed {
				taskText = completedStyle.Render(task.Label)
			}
			s += fmt.Sprintf("%s %s\n", cursor, taskText)
		}
		s += "\n⚔️  Press 'n' to add a new daily battle ⚔️."
	}
	return mainStyle.Render(s)
}

// TODO: Move into DB file
func saveTasks(tasks []task) error {
	data, err := json.Marshal(tasks)
	if err != nil {
		return err
	}
	return os.WriteFile("tasks.json", data, 0644)
}

func loadTasks() ([]task, error) {
	data, err := os.ReadFile("tasks.json")
	if err != nil {
		return nil, err // File might not exist on first run; handle accordingly.
	}
	var tasks []task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}
