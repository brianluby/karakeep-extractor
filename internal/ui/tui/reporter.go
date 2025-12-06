package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// BubbleTeaReporter adapts domain.ProgressReporter events to Bubble Tea messages.
type BubbleTeaReporter struct {
	program *tea.Program
}

func NewBubbleTeaReporter(p *tea.Program) *BubbleTeaReporter {
	return &BubbleTeaReporter{program: p}
}

func (r *BubbleTeaReporter) Start(total int, title string) {
	r.program.Send(MsgStart{Total: total, Title: title})
}

func (r *BubbleTeaReporter) Increment() {
	r.program.Send(MsgProgress{Increment: 1})
}

func (r *BubbleTeaReporter) SetStatus(status string) {
	r.program.Send(MsgStatus{Status: status})
}

func (r *BubbleTeaReporter) Log(message string) {
	r.program.Send(MsgLog{Message: message})
}

func (r *BubbleTeaReporter) Error(err error) {
	r.program.Send(MsgError{Err: err})
}

func (r *BubbleTeaReporter) Finish(summary string) {
	r.program.Send(MsgDone{Summary: summary})
}
