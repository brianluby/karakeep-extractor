package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LogModel displays a tail of recent messages.
type LogModel struct {
	messages []string
	maxLines int
}

func NewLogModel() LogModel {
	return LogModel{
		messages: []string{},
		maxLines: 5, // Keep last 5 lines
	}
}

func (m LogModel) Init() tea.Cmd {
	return nil
}

func (m LogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MsgLog:
		m.messages = append(m.messages, fmt.Sprintf("[INFO] %s", msg.Message))
		if len(m.messages) > m.maxLines {
			m.messages = m.messages[1:]
		}
	case MsgError:
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		m.messages = append(m.messages, style.Render(fmt.Sprintf("[ERROR] %v", msg.Err)))
		if len(m.messages) > m.maxLines {
			m.messages = m.messages[1:]
		}
	}
	return m, nil
}

func (m LogModel) View() string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")),
		Padding(0, 1).
		Width(60)

	content := strings.Join(m.messages, "\n")
	if len(m.messages) == 0 {
		content = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("No logs yet...")
	}
	
	return style.Render(content)
}
