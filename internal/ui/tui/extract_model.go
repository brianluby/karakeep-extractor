package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ExtractModel struct {
	spinner spinner.Model
	count   int
	status  string
}

func NewExtractModel() ExtractModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return ExtractModel{
		spinner: s,
		count:   0,
		status:  "Initializing...",
	}
}

func (m ExtractModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m ExtractModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case MsgStart:
		m.count = 0
		m.status = "Starting..."
	
	case MsgProgress:
		m.count += msg.Increment
	
	case MsgStatus:
		m.status = msg.Status

	case spinner.TickMsg:
		var newSpinner spinner.Model
		newSpinner, cmd = m.spinner.Update(msg)
		m.spinner = newSpinner
	}

	return m, cmd
}

func (m ExtractModel) View() string {
	return fmt.Sprintf("\n %s %s\n\n Processed: %d\n", m.spinner.View(), m.status, m.count)
}
