package service

import (
	"context"
	"fmt"
	"io"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
	"github.com/brianluby/karakeep-extractor/internal/ui"
)

type Ranker struct {
	repo domain.RankingRepository
}

func NewRanker(repo domain.RankingRepository) *Ranker {
	return &Ranker{repo: repo}
}

func (r *Ranker) Rank(ctx context.Context, limit int, sortBy string, output io.Writer) error {
	var sortOption domain.RankSortOption
	switch sortBy {
	case "stars":
		sortOption = domain.SortByStars
	case "forks":
		sortOption = domain.SortByForks
	case "updated":
		sortOption = domain.SortByUpdated
	default:
		return fmt.Errorf("invalid sort option: %s (valid: stars, forks, updated)", sortBy)
	}

	repos, err := r.repo.GetRankedRepos(ctx, limit, sortOption)
	if err != nil {
		return fmt.Errorf("failed to get ranked repos: %w", err)
	}

	if len(repos) == 0 {
		fmt.Fprintln(output, "No repositories found.")
		return nil
	}

	// Use Pager if output is stdout (which it usually is from main)
	// But here we are passed an io.Writer. If it's not stdout, UsePager fallback handles it.
	
	// If we want to use the pager logic:
	return ui.UsePager(output, func(w io.Writer) error {
		renderer := ui.NewTableRenderer(w)
		return renderer.Render(repos)
	})
}
