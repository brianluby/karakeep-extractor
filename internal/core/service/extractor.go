package service

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

// Extractor orchestrates the bookmark fetching, filtering, and saving process.
type Extractor struct {
	Source     domain.BookmarkSource
	Repository domain.RepoRepository
}

// NewExtractor creates a new Extractor service.
func NewExtractor(source domain.BookmarkSource, repository domain.RepoRepository) *Extractor {
	return &Extractor{
		Source:     source,
		Repository: repository,
	}
}

// Extract fetches bookmarks, filters for GitHub repos, normalizes URLs, and saves them.
func (e *Extractor) Extract(ctx context.Context, reporter domain.ProgressReporter) error {
	reporter.Start(-1, "Extracting bookmarks")
	currentPage := 1
	extractedCount := 0

	for {
		reporter.SetStatus(fmt.Sprintf("Fetching page %d", currentPage))
		bookmarks, err := e.Source.FetchBookmarks(ctx, currentPage)
		if err != nil {
			reporter.Error(err)
			return fmt.Errorf("failed to fetch bookmarks page %d: %w", currentPage, err)
		}

		if len(bookmarks) == 0 {
			break // No more bookmarks
		}

		for _, bm := range bookmarks {
			targetURL := bm.Content.URL
			normalizedRepoID, isGitHub := NormalizeGitHubURL(targetURL)
			if !isGitHub {
				continue // Skip non-GitHub URLs
			}

			// Check for malformed URL after normalization attempt
			if normalizedRepoID == "" {
				reporter.Log(fmt.Sprintf("Skipping malformed URL in bookmark ID %s: %s", bm.ID, targetURL))
				continue
			}

			exists, err := e.Repository.Exists(ctx, normalizedRepoID)
			if err != nil {
				reporter.Log(fmt.Sprintf("Error checking existence for %s: %v", normalizedRepoID, err))
				continue
			}
			if exists {
				// log.Printf("Skipping duplicate repo: %s", normalizedRepoID)
				continue
			}

			// Determine Title
			title := bm.Content.Title
			if bm.Title != nil && *bm.Title != "" {
				title = *bm.Title
			}

			repo := domain.ExtractedRepo{
				RepoID:   normalizedRepoID,
				URL:      targetURL, // Keep original URL for now, can be normalized later if needed
				SourceID: bm.ID,
				Title:    title,
				FoundAt:  time.Now(),
			}

			if err := e.Repository.Save(ctx, repo); err != nil {
				reporter.Log(fmt.Sprintf("Error saving repo %s: %v", normalizedRepoID, err))
				continue
			}
			extractedCount++
			reporter.Increment()
		}
		currentPage++
	}
	reporter.Finish(fmt.Sprintf("Extracted %d new repositories", extractedCount))
	return nil
}

var githubDomainRegex = regexp.MustCompile(`^(www\.)?github\.com$`)
var repoPathRegex = regexp.MustCompile(`^/?([^/]+)/([^/]+)`) // Matches /owner/repo

// NormalizeGitHubURL attempts to normalize a GitHub URL to "owner/repo" format.
// Returns the normalized string and a boolean indicating if it's a GitHub URL with owner/repo.
func NormalizeGitHubURL(rawURL string) (string, bool) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", false // Malformed URL
	}

	if !githubDomainRegex.MatchString(strings.ToLower(u.Host)) {
		return "", false // Not a GitHub domain
	}

	matches := repoPathRegex.FindStringSubmatch(u.Path)
	if len(matches) < 3 {
		return "", false // Path doesn't contain owner/repo
	}

	owner := matches[1]
	repo := strings.TrimSuffix(matches[2], ".git") // Remove .git suffix if present

	return fmt.Sprintf("%s/%s", owner, repo), true
}
