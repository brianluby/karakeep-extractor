package tui

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
)

// Verify BubbleTeaReporter implements ProgressReporter via msg verification
func TestBubbleTeaReporter_SendMsgs(t *testing.T) {
	program := tea.NewProgram(nil) // Dummy program, won't run
	reporter := NewBubbleTeaReporter(program)

	// This test is slightly tricky because we need to capture the messages sent to the program.
	// Bubble Tea doesn't easily expose a way to spy on Send() without mocking the Program struct,
	// which is concrete.
	// However, we can test that the methods don't panic and likely compile correctly.
	// A true integration test would be better here.

	// For now, just simple coverage that calls work
	reporter.Start(10, "Test Task")
	reporter.SetStatus("Processing")
	reporter.Increment()
	reporter.Log("Info log")
	reporter.Finish("Done")
}
