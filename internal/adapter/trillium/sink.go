package trillium

import (
	"context"
	"fmt"
	"time"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
	"github.com/brianluby/karakeep-extractor/internal/ui"
)

// TrilliumSink sends repositories to Trillium as a note.
type TrilliumSink struct {
	client *TrilliumClient
}

func NewSink(client *TrilliumClient) *TrilliumSink {
	return &TrilliumSink{client: client}
}

func (s *TrilliumSink) Send(ctx context.Context, repos []domain.ExtractedRepo) error {
	// 1. Format as Markdown
	formatter := ui.NewMarkdownFormatter()
	content := formatter.FormatTable(repos)

	// 2. Create Note
	title := fmt.Sprintf("GitHub Rankings - %s", time.Now().Format("2006-01-02 15:04:05"))
	
	return s.client.CreateNote(ctx, title, content)
}
