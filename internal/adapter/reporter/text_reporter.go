package reporter

import (
	"log"
)

// TextReporter is the standard stdout/stderr implementation of ProgressReporter.
type TextReporter struct {
	// We can add fields here if we need to track count manually for summary
	count int
}

func NewTextReporter() *TextReporter {
	return &TextReporter{}
}

func (r *TextReporter) Start(total int, title string) {
	log.Printf("Starting: %s (Total: %d)", title, total)
}

func (r *TextReporter) Increment() {
	r.count++
	// Optional: log every N items if desired, but for now we keep it quiet to match legacy or just log completion
}

func (r *TextReporter) SetStatus(status string) {
	// In text mode, setting status usually means logging what we are working on
	log.Println(status)
}

func (r *TextReporter) Log(message string) {
	log.Println(message)
}

func (r *TextReporter) Error(err error) {
	log.Printf("Error: %v", err)
}

func (r *TextReporter) Finish(summary string) {
	log.Printf("Finished: %s", summary)
}
