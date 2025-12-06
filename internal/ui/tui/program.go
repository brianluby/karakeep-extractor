package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

// Run starts the Bubble Tea program.
// task: A closure/function that performs the actual work and uses the internal reporter.
func Run(ctx context.Context, mode string, task func(domain.ProgressReporter) error) error {
	var opMode OperationMode
	if mode == "enrich" {
		opMode = ModeEnrich
	} else {
		opMode = ModeExtract
	}

	// Use WithAltScreen to restore terminal on exit
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Create a reporter bound to this program
	reporter := NewBubbleTeaReporter(p)

	// The work function needs to be executed.
	// In UI-Driven model (Question 3/Option A in Spec), the Program owns the loop.
	// But we need to pass the reporter to the work function.
	// The RootModel needs to know how to start the work.
	
	// We need to inject the bound reporter into the work function stored in the model.
	// But RootModel is a value type in Update, so we can't easily mutating it after Init?
	// Actually, Init() returns a Cmd. We can use that Cmd to start a goroutine that runs `task(reporter)`.
	
	// Let's redefine RootModel's Work field to accept the reporter we just created?
	// Or better: The `MsgWorkerStart` handler in Update() should trigger the work.
	// But Update() is pure. It can return a Cmd.
	
	// Correct approach:
	// 1. Pass the PRE-BOUND worker closure to NewRootModel?
	//    No, we only have 'task func(Reporter)'. We need 'func() error'.
	// 2. Wrap it here.
	
	boundWork := func() error {
		return task(reporter)
	}
	
	// We need to pass this boundWork to the model so it can call it in a Cmd.
	// But `tea.Cmd` expects `func() tea.Msg`.
	
	// Let's update RootModel to accept a tea.Cmd-compatible starter or just the closure.
	// We will modify NewRootModel to accept `func() error`.
	
	// Wait, I can't modify RootModel easily from here if I already defined it in model.go with a different signature.
	// Let's adjust model.go first to align with this cleaner plan.
	
	// Actually, let's look at `internal/ui/tui/model.go`.
	// It has `Work func(domain.ProgressReporter) error`.
	// And Init() returns `MsgWorkerStart`.
	// And Update() handles `MsgWorkerStart`.
	
	// But inside Update(), we don't have the `reporter` instance we created here!
	// The reporter needs `p` (the Program).
	// This is a circular dependency if we try to put the reporter INSIDE the model.
	
	// Standard Bubble Tea pattern for external events:
	// Run the worker in a separate goroutine OUTSIDE the model, and have it send messages to `p`.
	// We don't need the model to "own" the worker execution via Update/Cmd if the worker is purely driving UI via messages.
	
	// So:
	go func() {
		// Wait for program to be ready? 
		// Usually fine to start sending messages immediately, they buffer.
		if err := boundWork(); err != nil {
			p.Send(MsgFatal{Err: err})
		} else {
			p.Send(MsgDone{Summary: "Completed successfully."})
		}
	}()
	
	_, err := p.Run()
	return err
}
