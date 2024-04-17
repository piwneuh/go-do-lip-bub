package main

import (
	"fmt"

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
	label     string
	completed bool
}

func main() {
	initialTasks := []task{
		{label: "Wake up sleeepy head 🛏️ 💤", completed: false},
		{label: "Brush your teeth 🦷", completed: false},
		{label: "Get some ☕", completed: false},
	}
	p := tea.NewProgram(model{tasks: initialTasks})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

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
				m.tasks[m.selected].completed = !m.tasks[m.selected].completed
			} else if m.input != "" {
				m.tasks = append(m.tasks, task{label: m.input, completed: false})
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

func (m model) View() string {
	var s string
	if m.creating {
		s += inputStyle.Render(fmt.Sprintf("Add Task: %s", m.input))
	} else {
		s += "🧙 Mornin', adventurer! ☁️⚡\n"
		for i, task := range m.tasks {
			cursor := " "
			if i == m.selected {
				cursor = pointerStyle.Render(">")
			}
			taskText := taskStyle.Render(task.label)
			if task.completed {
				taskText = completedStyle.Render(task.label)
			}
			s += fmt.Sprintf("%s %s\n", cursor, taskText)
		}
		s += "\nPress 'n' to add a new task."
	}
	return mainStyle.Render(s)
}
