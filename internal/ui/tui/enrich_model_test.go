package tui

import (
	"testing"
)

func TestEnrichModel_Update_Stats(t *testing.T) {
	m := NewEnrichModel()

	// Simulate Success
	newM, _ := m.Update(MsgSuccess{})
	m = newM.(EnrichModel)
	if m.stats.SuccessCount != 1 {
		t.Errorf("Expected SuccessCount 1, got %d", m.stats.SuccessCount)
	}

	// Simulate Failure
	newM, _ = m.Update(MsgFailure{})
	m = newM.(EnrichModel)
	if m.stats.FailureCount != 1 {
		t.Errorf("Expected FailureCount 1, got %d", m.stats.FailureCount)
	}

	// Simulate Skipped
	newM, _ = m.Update(MsgSkipped{})
	m = newM.(EnrichModel)
	if m.stats.SkippedCount != 1 {
		t.Errorf("Expected SkippedCount 1, got %d", m.stats.SkippedCount)
	}
}

func TestEnrichModel_Update_Progress(t *testing.T) {
	m := NewEnrichModel()
	m.total = 10

	// MsgProgress should increment current count
	newM, _ := m.Update(MsgProgress{Increment: 1})
	m = newM.(EnrichModel)
	if m.current != 1 {
		t.Errorf("Expected current 1, got %d", m.current)
	}
}
