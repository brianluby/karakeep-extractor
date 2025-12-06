package tui

import (
	"io"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type testModel struct {
	msgs chan tea.Msg
}

func (m testModel) Init() tea.Cmd { return nil }

func (m testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Non-blocking send to channel for verification
	select {
	case m.msgs <- msg:
	default:
	}

	if _, ok := msg.(MsgDone); ok {
		return m, tea.Quit
	}
	return m, nil
}

func (m testModel) View() string { return "" }

// Verify BubbleTeaReporter sends messages correctly
func TestBubbleTeaReporter_SendMsgs(t *testing.T) {
	msgCh := make(chan tea.Msg, 10)
	model := testModel{msgs: msgCh}
	
	// Initialize program with no IO to avoid interfering with test output
	program := tea.NewProgram(model, tea.WithInput(nil), tea.WithOutput(io.Discard))
	reporter := NewBubbleTeaReporter(program)

	// Run program in a separate goroutine so it can process messages
	go func() {
		program.Run()
	}()

	// Send messages
	reporter.Start(10, "Test Task")
	reporter.SetStatus("Processing")
	reporter.Increment()
	reporter.Log("Info log")
	reporter.Error(nil)
	reporter.Finish("Done")

	// Verify we receive at least some messages and eventually quit
	// We expect: MsgStart, MsgStatus, MsgProgress, MsgLog, MsgError, MsgDone
	expectedCount := 6
	receivedCount := 0

	for i := 0; i < expectedCount; i++ {
		select {
		case <-msgCh:
			receivedCount++
		}
	}

	if receivedCount != expectedCount {
		t.Errorf("Expected %d messages, got %d", expectedCount, receivedCount)
	}
}
