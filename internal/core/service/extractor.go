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
	reporter.SetStatus("Fetching all bookmarks...")
	bookmarks, err := e.Source.FetchBookmarks(ctx)
	if err != nil {
		reporter.Error(err)
		return fmt.Errorf("failed to fetch all bookmarks: %w", err)
	}

	if len(bookmarks) == 0 {
		reporter.Log("No bookmarks found.")
		reporter.Finish("Extraction complete: 0 new repositories.")
		return nil
	}
	
	reporter.Start(len(bookmarks), "Processing bookmarks")
	extractedCount := 0
	
	// Regex to find potential links in text (simplified)
	linkRegex := regexp.MustCompile(`https?://github\.com/[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+`)

	for _, bm := range bookmarks {
		// Candidate URLs: Main URL + any found in HTML content
		candidates := []string{bm.Content.URL}
		
		if bm.Content.HTMLContent != "" {
			matches := linkRegex.FindAllString(bm.Content.HTMLContent, -1)
			candidates = append(candidates, matches...)
		}

		// Deduplicate candidates for this bookmark to avoid processing same repo twice
		uniqueRepos := make(map[string]string) // normalizedID -> originalURL

		foundNew := false
		for _, rawURL := range candidates {
			normalizedRepoID, isGitHub := NormalizeGitHubURL(rawURL)
			if !isGitHub || normalizedRepoID == "" {
				continue
			}
			uniqueRepos[normalizedRepoID] = rawURL
		}

		for normalizedRepoID, originalURL := range uniqueRepos {
			exists, err := e.Repository.Exists(ctx, normalizedRepoID)
			if err != nil {
				reporter.Log(fmt.Sprintf("Error checking existence for %s: %v", normalizedRepoID, err))
				continue
			}
			if exists {
				continue
			}

			// Determine Title (Use bookmark title, or fallback to repo ID if finding multiple?)
			title := bm.Content.Title
			if bm.Title != nil && *bm.Title != "" {
				title = *bm.Title
			}

			repo := domain.ExtractedRepo{
				RepoID:   normalizedRepoID,
				URL:      originalURL,
				SourceID: bm.ID,
				Title:    title,
				FoundAt:  time.Now(),
			}

			if err := e.Repository.Save(ctx, repo); err != nil {
				reporter.Log(fmt.Sprintf("Error saving repo %s: %v", normalizedRepoID, err))
				reporter.RecordFailure()
				continue
			}
			extractedCount++
			foundNew = true
		}
		
		if foundNew {
			reporter.RecordSuccess() // Treat "Processed & Found Repo" as Success
		} else {
			reporter.RecordSkipped() // Treat "Processed & No New Repo" as Skipped
		}
		reporter.Increment()
	}
	reporter.Finish(fmt.Sprintf("Extraction complete: %d new repositories found.", extractedCount))
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

	owner := strings.ToLower(matches[1])
	repo := strings.TrimSuffix(matches[2], ".git") // Remove .git suffix if present

	// Blacklist of reserved paths that look like owner/repo but aren't
	reservedPaths := map[string]bool{
		"marketplace":      true,
		"apps":             true,
		"sponsors":         true,
		"advisories":       true,
		"topics":           true,
		"search":           true,
		"login":            true,
		"join":             true,
		"features":         true,
		"pricing":          true,
		"enterprise":       true,
		"customer-stories": true,
		"security":         true,
		"collections":      true,
		"new":              true,
		"settings":         true,
		"site":             true,
		"about":            true,
		"contact":          true,
		"organizations":    true,
	}

	if reservedPaths[owner] {
		return "", false
	}

	return fmt.Sprintf("%s/%s", matches[1], repo), true // Return original casing for owner
}