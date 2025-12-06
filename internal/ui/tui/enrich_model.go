package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EnrichModel struct {
	progress progress.Model
	total    int
	current  int
	status   string
}

func NewEnrichModel() EnrichModel {
	return EnrichModel{
		progress: progress.New(progress.WithDefaultGradient()),
		total:    0,
		current:  0,
		status:   "Waiting...",
	}
}

func (m EnrichModel) Init() tea.Cmd {
	return nil
}

func (m EnrichModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case MsgStart:
		m.total = msg.Total
		m.current = 0
		m.status = "Starting..."
	
	case MsgProgress:
		m.current += msg.Increment
		if m.total > 0 {
			pct := float64(m.current) / float64(m.total)
			cmd = m.progress.SetPercent(pct)
		}
	
	case MsgStatus:
		m.status = msg.Status

	case progress.FrameMsg:
		newModel, newCmd := m.progress.Update(msg)
		m.progress = newModel.(progress.Model)
		cmd = newCmd
	}

	return m, cmd
}

func (m EnrichModel) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(m.status) + "\n"
}

const padding = 2
