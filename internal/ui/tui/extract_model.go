package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ExtractModel struct {
	spinner spinner.Model
	count   int
	status  string
	stats   ProgressStats
}

func NewExtractModel() ExtractModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return ExtractModel{
		spinner: s,
		count:   0,
		status:  "Initializing...",
		stats:   ProgressStats{},
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
		m.stats = ProgressStats{}
	
	case MsgProgress:
		m.count += msg.Increment
	
	case MsgStatus:
		m.status = msg.Status

	case MsgSuccess:
		m.stats.SuccessCount++
	
	case MsgFailure:
		m.stats.FailureCount++
	
	case MsgSkipped:
		m.stats.SkippedCount++

	case spinner.TickMsg:
		var newSpinner spinner.Model
		newSpinner, cmd = m.spinner.Update(msg)
		m.spinner = newSpinner
	}

	return m, cmd
}

func (m ExtractModel) View() string {
	statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	statsStr := fmt.Sprintf("Processed: %d | Failed: %d | Skipped: %d", m.stats.SuccessCount, m.stats.FailureCount, m.stats.SkippedCount)
	pad := strings.Repeat(" ", 2)

	return fmt.Sprintf("\n%s %s %s\n\n%s%s\n", pad, m.spinner.View(), m.status, pad, statsStyle.Render(statsStr))
}