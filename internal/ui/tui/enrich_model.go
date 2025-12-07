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
	stats    ProgressStats
}

func NewEnrichModel() EnrichModel {
	return EnrichModel{
		progress: progress.New(progress.WithDefaultGradient()),
		total:    0,
		current:  0,
		status:   "Waiting...",
		stats:    ProgressStats{},
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
		m.stats = ProgressStats{} // Reset stats
	
	case MsgProgress:
		m.current += msg.Increment
		if m.total > 0 {
			pct := float64(m.current) / float64(m.total)
			cmd = m.progress.SetPercent(pct)
		}
	
	case MsgStatus:
		m.status = msg.Status

	case MsgSuccess:
		m.stats.SuccessCount++
	
	case MsgFailure:
		m.stats.FailureCount++
	
	case MsgSkipped:
		m.stats.SkippedCount++

	case progress.FrameMsg:
		newModel, newCmd := m.progress.Update(msg)
		m.progress = newModel.(progress.Model)
		cmd = newCmd
	}

	return m, cmd
}

func (m EnrichModel) View() string {
	pad := strings.Repeat(" ", padding)
	
	statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	statsStr := fmt.Sprintf("Success: %d | Failed: %d | Skipped: %d", m.stats.SuccessCount, m.stats.FailureCount, m.stats.SkippedCount)

	return "\n" +
		pad + m.progress.View() + "\n" +
		pad + statsStyle.Render(statsStr) + "\n\n" +
		pad + lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(m.status) + "\n"
}

const padding = 2