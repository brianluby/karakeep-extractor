package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type OperationMode int

const (
	ModeIdle OperationMode = iota
	ModeEnrich
	ModeExtract
)

type AppState int

const (
	StateIdle AppState = iota
	StateRunning
	StateDone
	StateError
	StateFatal
)

// RootModel is the top-level Bubble Tea model.
type RootModel struct {
	State     AppState
	Mode      OperationMode
	TaskTitle string
	Summary   string
	FatalErr  error

	// Child Models
	EnrichModel EnrichModel
	ExtractModel ExtractModel
	LogModel    LogModel
}

func NewRootModel(mode OperationMode) RootModel {
	m := RootModel{
		State:        StateIdle,
		Mode:         mode,
		EnrichModel:  NewEnrichModel(),
		ExtractModel: NewExtractModel(),
		LogModel:     NewLogModel(),
	}
	return m
}

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(
		m.EnrichModel.Init(),
		m.ExtractModel.Init(),
	)
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		// We can propagate this to children if they need to resize
		// For now, we just acknowledge it. If we used viewport, we'd pass it on.
		// progress bar might resize automatically or need setting width.
		if m.Mode == ModeEnrich {
			m.EnrichModel.progress.Width = msg.Width - padding*2 - 4
			if m.EnrichModel.progress.Width > 80 {
				m.EnrichModel.progress.Width = 80
			}
		}

	case MsgStart:
		m.State = StateRunning
		m.TaskTitle = msg.Title
		// Delegate
		if m.Mode == ModeEnrich {
			newM, cmd := m.EnrichModel.Update(msg)
			m.EnrichModel = newM.(EnrichModel)
			cmds = append(cmds, cmd)
		} else {
			newM, cmd := m.ExtractModel.Update(msg)
			m.ExtractModel = newM.(ExtractModel)
			cmds = append(cmds, cmd)
		}

	case MsgProgress:
		if m.Mode == ModeEnrich {
			newM, cmd := m.EnrichModel.Update(msg)
			m.EnrichModel = newM.(EnrichModel)
			cmds = append(cmds, cmd)
		} else {
			newM, cmd := m.ExtractModel.Update(msg)
			m.ExtractModel = newM.(ExtractModel)
			cmds = append(cmds, cmd)
		}

	case MsgStatus:
		// Update status bar or current item
		if m.Mode == ModeEnrich {
			newM, cmd := m.EnrichModel.Update(msg)
			m.EnrichModel = newM.(EnrichModel)
			cmds = append(cmds, cmd)
		} else {
			newM, cmd := m.ExtractModel.Update(msg)
			m.ExtractModel = newM.(ExtractModel)
			cmds = append(cmds, cmd)
		}

	case MsgLog:
		newLog, cmd := m.LogModel.Update(msg)
		m.LogModel = newLog.(LogModel)
		cmds = append(cmds, cmd)

	case MsgError:
		newLog, cmd := m.LogModel.Update(msg)
		m.LogModel = newLog.(LogModel)
		cmds = append(cmds, cmd)

	case MsgDone:
		m.State = StateDone
		m.Summary = msg.Summary
		// return m, tea.Quit // Don't quit yet, let user see summary
		return m, nil

	case MsgFatal:
		m.State = StateFatal
		m.FatalErr = msg.Err
		return m, tea.Quit // Fatal errors can quit or stay? Let's quit for now or wait for key?
		// Spec says "The system MUST restore the terminal... upon exit or failure".
		// If we Quit, we restore.
	
	default:
		// Forward specific messages to child models if needed
		if m.Mode == ModeEnrich {
			newM, cmd := m.EnrichModel.Update(msg)
			m.EnrichModel = newM.(EnrichModel)
			cmds = append(cmds, cmd)
		} else {
			newM, cmd := m.ExtractModel.Update(msg)
			m.ExtractModel = newM.(ExtractModel)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	if m.FatalErr != nil {
		return fmt.Sprintf("Fatal Error: %v\n", m.FatalErr)
	}
	
	if m.State == StateDone {
		return fmt.Sprintf("\n%s\n\nDone! %s\nPress Ctrl+C or q to quit.\n", m.LogModel.View(), m.Summary)
	}

	var s strings.Builder

	s.WriteString(fmt.Sprintf("Karakeep Extractor: %s\n\n", m.TaskTitle))

	if m.Mode == ModeEnrich {
		s.WriteString(m.EnrichModel.View())
	} else {
		s.WriteString(m.ExtractModel.View())
	}

	s.WriteString("\n\n")
	s.WriteString(m.LogModel.View())
	s.WriteString("\nPress Ctrl+C to quit.\n")

	return s.String()
}

// Messages
type MsgStart struct { Total int; Title string }
type MsgProgress struct { Increment int }
type MsgStatus struct { Status string }
type MsgLog struct { Message string }
type MsgError struct { Err error }
type MsgDone struct { Summary string }
type MsgFatal struct { Err error }

