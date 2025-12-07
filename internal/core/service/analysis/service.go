package analysis

import (
	"context"
	"fmt"
	"strings"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

type LLMProvider interface {
	SendMessage(ctx context.Context, req domain.AnalysisRequest) (string, error)
}

type Service struct {
	repo domain.RankingRepository
	llm  LLMProvider
}

func NewService(repo domain.RankingRepository, llm LLMProvider) *Service {
	return &Service{
		repo: repo,
		llm:  llm,
	}
}

func (s *Service) Analyze(ctx context.Context, query string, limit int, langFilter string, tagFilter string, minStars int, maxStars int) (string, error) {
	// Fetch a large batch to allow for local filtering
	// If the user asks for a range like 500-1000, and we only fetch top 500 by stars, we might miss them if they are further down.
	// We might need to fetch MORE if a specific range is requested that isn't at the top.
	// For now, let's bump the fetch limit if filtering by stars to ensure we catch them.
	fetchLimit := 2000 
	if limit > fetchLimit {
		fetchLimit = limit
	}

	// Default sort by stars to get "best" repos first
	repos, err := s.repo.GetRankedRepos(ctx, fetchLimit, domain.SortByStars, tagFilter)
	if err != nil {
		return "", fmt.Errorf("failed to fetch repos: %w", err)
	}

	// Filter
	var filtered []domain.ExtractedRepo
	for _, r := range repos {
		// Language Filter
		if langFilter != "" {
			if r.Language == nil || !strings.EqualFold(*r.Language, langFilter) {
				continue
			}
		}

		// Stars Filter
		stars := 0
		if r.Stars != nil {
			stars = *r.Stars
		}
		if stars < minStars {
			continue
		}
		if maxStars > 0 && stars > maxStars {
			continue
		}

		filtered = append(filtered, r)
	}

	// Apply final limit
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}

	if len(filtered) == 0 {
		return "No repositories found matching your criteria.", nil
	}

	// Build Prompt
	msgs, err := BuildMessages(query, filtered)
	if err != nil {
		return "", err
	}

	// Call LLM
	req := domain.AnalysisRequest{
		Messages: msgs,
	}
	return s.llm.SendMessage(ctx, req)
}
