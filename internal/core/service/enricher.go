package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

// Enricher orchestrates the enrichment process.
type Enricher struct {
	repo   domain.RepoRepository
	client domain.GitHubClient
}

func NewEnricher(repo domain.RepoRepository, client domain.GitHubClient) *Enricher {
	return &Enricher{
		repo:   repo,
		client: client,
	}
}

type EnrichmentResult struct {
	RepoID string
	Status domain.EnrichmentStatus
	Err    error
}

func (e *Enricher) EnrichBatch(ctx context.Context, limit int, force bool, workers int, reporter domain.ProgressReporter) (int, int, error) {
	repos, err := e.repo.GetReposForEnrichment(ctx, limit, force)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get repos for enrichment: %w", err)
	}

	if len(repos) == 0 {
		return 0, 0, nil
	}

	// Initialize Reporter
	reporter.Start(len(repos), "Enriching repositories")

	// Create a cancellable context for this batch
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobCh := make(chan *domain.ExtractedRepo, len(repos))
	resCh := make(chan EnrichmentResult, len(repos))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for repo := range jobCh {
				// Check context cancellation to exit early
				if ctx.Err() != nil {
					return
				}
				e.processRepo(ctx, repo, resCh, reporter)
			}
		}()
	}

	// Enqueue jobs
	for _, repo := range repos {
		jobCh <- repo
	}
	close(jobCh)

	// Wait for workers in a separate goroutine to close result channel
	go func() {
		wg.Wait()
		close(resCh)
	}()

	successCount := 0
	notFoundCount := 0
	errCount := 0

	for res := range resCh {
		switch res.Status {
		case domain.StatusSuccess:
			successCount++
		case domain.StatusNotFound:
			notFoundCount++
		default:
			errCount++
			// If Fail Fast is triggered (Rate Limit), we should probably return early.
			// Ideally, processRepo returns a specific error type we can check.
			if res.Status == domain.StatusAPIError && res.Err != nil && errors.Is(res.Err, domain.ErrRateLimitExceeded) {
				cancel() // Stop all other workers immediately
				reporter.Error(domain.ErrRateLimitExceeded)
				return successCount, errCount, domain.ErrRateLimitExceeded
			}
		}
		reporter.Increment()
	}
	
	reporter.Finish(fmt.Sprintf("Enriched: %d, Not Found: %d, Failed: %d", successCount, notFoundCount, errCount))

	return successCount, errCount, nil
}

func (e *Enricher) processRepo(ctx context.Context, repo *domain.ExtractedRepo, resCh chan<- EnrichmentResult, reporter domain.ProgressReporter) {
	// Parse Owner/Repo from URL or RepoID
	// Assuming RepoID is already "owner/repo" as per domain
	parts := strings.Split(repo.RepoID, "/")
	if len(parts) != 2 {
		err := fmt.Errorf("invalid repo id format")
		reporter.Log(fmt.Sprintf("Skipping %s: %v", repo.RepoID, err))
		resCh <- EnrichmentResult{RepoID: repo.RepoID, Status: domain.StatusAPIError, Err: err}
		return
	}
	owner, name := parts[0], parts[1]

	reporter.SetStatus(fmt.Sprintf("Enriching %s", repo.RepoID))
	stats, _, err := e.client.GetRepoStats(ctx, owner, name)
	
	update := domain.RepoEnrichmentUpdate{
		RepoID: repo.RepoID,
	}

	if err != nil {
		if errors.Is(err, domain.ErrRateLimitExceeded) {
			// Critical error, handled by orchestrator to stop
			update.EnrichmentStatus = domain.StatusAPIError // Or keep pending?
			// We don't update DB on rate limit to retry later? Or mark API_ERROR?
			// Plan says "Save progress and exit". API_ERROR allows retry if logic permits.
			reporter.Log(fmt.Sprintf("Rate limit exceeded for %s", repo.RepoID))
			resCh <- EnrichmentResult{RepoID: repo.RepoID, Status: domain.StatusAPIError, Err: err}
			return
		}
		if errors.Is(err, domain.ErrRepoNotFound) {
			update.EnrichmentStatus = domain.StatusNotFound
			reporter.Log(fmt.Sprintf("Repo not found: %s", repo.RepoID))
		} else {
			update.EnrichmentStatus = domain.StatusAPIError
			reporter.Log(fmt.Sprintf("API error for %s: %v", repo.RepoID, err))
		}
	} else {
		update.Stats = stats
		update.EnrichmentStatus = domain.StatusSuccess
	}

	// Persist
	if saveErr := e.repo.UpdateRepoEnrichment(ctx, update); saveErr != nil {
		reporter.Log(fmt.Sprintf("Save failed for %s: %v", repo.RepoID, saveErr))
		resCh <- EnrichmentResult{RepoID: repo.RepoID, Status: domain.StatusAPIError, Err: fmt.Errorf("save failed: %w", saveErr)}
		return
	}

	resCh <- EnrichmentResult{RepoID: repo.RepoID, Status: update.EnrichmentStatus, Err: err}
}
