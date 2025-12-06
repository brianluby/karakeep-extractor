package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
	"github.com/brianluby/karakeep-extractor/internal/core/service"
)

// MockBookmarkSource for testing extractor service
type mockBookmarkSource struct {
	bookmarks [][]domain.RawBookmark
	callCount int
}

func (m *mockBookmarkSource) FetchBookmarks(ctx context.Context) ([]domain.RawBookmark, error) {
	if m.callCount >= len(m.bookmarks) {
		return nil, nil // No more pages (shouldn't happen with non-paginated API)
	}
	res := m.bookmarks[m.callCount]
	m.callCount++ // Still increment to "return" new content if multiple calls are made
	return res, nil
}

// MockRepoRepository for testing extractor service
type mockRepoRepository struct {
	repos map[string]domain.ExtractedRepo
}

func newMockRepoRepository() *mockRepoRepository {
	return &mockRepoRepository{
		repos: make(map[string]domain.ExtractedRepo),
	}
}

func (m *mockRepoRepository) Save(ctx context.Context, repo domain.ExtractedRepo) error {
	if _, exists := m.repos[repo.RepoID]; exists {
		return fmt.Errorf("duplicate RepoID: %s", repo.RepoID) // Simulate unique constraint
	}
	m.repos[repo.RepoID] = repo
	return nil
}

func (m *mockRepoRepository) Exists(ctx context.Context, repoID string) (bool, error) {
	_, exists := m.repos[repoID]
	return exists, nil
}

func (m *mockRepoRepository) GetReposForEnrichment(ctx context.Context, limit int, force bool) ([]*domain.ExtractedRepo, error) {
	return nil, nil
}

func (m *mockRepoRepository) UpdateRepoEnrichment(ctx context.Context, update domain.RepoEnrichmentUpdate) error {
	return nil
}

// MockReporter for testing
type mockReporter struct{}

func (m *mockReporter) Start(total int, title string)   {}
func (m *mockReporter) Increment()                      {}
func (m *mockReporter) SetStatus(status string)         {}
func (m *mockReporter) Log(message string)              {}
func (m *mockReporter) Error(err error)                 {}
func (m *mockReporter) Finish(summary string)           {}

func TestExtractService_Extract(t *testing.T) {
	testCases := []struct {
		name              string
		rawBookmarks      [][]domain.RawBookmark
		expectedRepos     []domain.ExtractedRepo
		expectedError     bool
		expectedLogOutput string // For malformed URLs
	}{
		{
			name: "single page, mixed bookmarks, no duplicates",
			rawBookmarks: [][]domain.RawBookmark{
				{
					{ID: "1", Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
						HTMLContent string `json:"htmlContent"`
					}{URL: "https://github.com/owner1/repo1", Title: "Repo One", HTMLContent: ""}},
					{ID: "2", Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
						HTMLContent string `json:"htmlContent"`
					}{URL: "https://example.com/not-github", Title: "Non-GitHub", HTMLContent: ""}},
					{ID: "3", Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
						HTMLContent string `json:"htmlContent"`
					}{URL: "https://github.com/owner2/repo2.git", Title: "Repo Two", HTMLContent: ""}},
				},
			},
			expectedRepos: []domain.ExtractedRepo{
				{RepoID: "owner1/repo1", URL: "https://github.com/owner1/repo1", SourceID: "1", Title: "Repo One"},
				{RepoID: "owner2/repo2", URL: "https://github.com/owner2/repo2.git", SourceID: "3", Title: "Repo Two"},
			},
			expectedError: false,
		},
		{
			name: "multiple pages, github links, deduplication",
			rawBookmarks: [][]domain.RawBookmark{
				{
					{ID: "10", Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
						HTMLContent string `json:"htmlContent"`
					}{URL: "https://github.com/ownerA/repoA", Title: "Repo A", HTMLContent: ""}},
					{ID: "11", Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
						HTMLContent string `json:"htmlContent"`
					}{URL: "https://github.com/ownerB/repoB", Title: "Repo B", HTMLContent: ""}},
					{ID: "12", Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
						HTMLContent string `json:"htmlContent"`
					}{URL: "http://github.com/ownerA/repoA/", Title: "Repo A Dupe", HTMLContent: ""}}, // Duplicate
					{ID: "13", Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
						HTMLContent string `json:"htmlContent"`
					}{URL: "https://github.com/ownerC/repoC?ref=master", Title: "Repo C", HTMLContent: ""}},
				},
			},
			expectedRepos: []domain.ExtractedRepo{
				{RepoID: "ownerA/repoA", URL: "https://github.com/ownerA/repoA", SourceID: "10", Title: "Repo A"},
				{RepoID: "ownerB/repoB", URL: "https://github.com/ownerB/repoB", SourceID: "11", Title: "Repo B"},
				{RepoID: "ownerC/repoC", URL: "https://github.com/ownerC/repoC?ref=master", SourceID: "13", Title: "Repo C"},
			},
			expectedError: false,
		},
		{
			name: "malformed URL handling",
			rawBookmarks: [][]domain.RawBookmark{
				{
					{ID: "1", Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
						HTMLContent string `json:"htmlContent"`
					}{URL: "https://github.com/good/repo", Title: "Good Repo", HTMLContent: ""}},
					{ID: "2", Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
						HTMLContent string `json:"htmlContent"`
					}{URL: "ftp://bad-url.com", Title: "Bad URL", HTMLContent: ""}},
					{ID: "3", Content: struct {
						URL         string `json:"url"`
						Title       string `json:"title"`
						Description string `json:"description"`
						HTMLContent string `json:"htmlContent"`
					}{URL: "invalid-url", Title: "Another Bad URL", HTMLContent: ""}},
				},
			},
			expectedRepos: []domain.ExtractedRepo{
				{RepoID: "good/repo", URL: "https://github.com/good/repo", SourceID: "1", Title: "Good Repo"},
			},
			expectedError: false,
			expectedLogOutput: "Skipping malformed URL", // Check if this appears in logs
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSource := &mockBookmarkSource{bookmarks: tc.rawBookmarks}
			mockRepo := newMockRepoRepository()
			extractor := service.NewExtractor(mockSource, mockRepo)

			err := extractor.Extract(context.Background(), &mockReporter{})
			if tc.expectedError && err == nil {
				t.Errorf("Expected an error but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got %v", err)
			}

			// Verify extracted repositories
			if len(mockRepo.repos) != len(tc.expectedRepos) {
				t.Fatalf("Expected %d repos, got %d", len(tc.expectedRepos), len(mockRepo.repos))
			}

			for _, expected := range tc.expectedRepos {
				found, ok := mockRepo.repos[expected.RepoID]
				if !ok {
					t.Errorf("RepoID %s not found in extracted repos", expected.RepoID)
					continue
				}
				// Compare relevant fields, ignore FoundAt for tests
				if found.RepoID != expected.RepoID || found.URL != expected.URL || found.SourceID != expected.SourceID || found.Title != expected.Title {
					t.Errorf("Mismatch for RepoID %s. Expected %+v, Got %+v", expected.RepoID, expected, found)
				}
			}
		})
	}
}

func TestNormalizeGitHubURL(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		isGitHub bool
	}{
		{"https://github.com/owner/repo", "owner/repo", true},
		{"http://github.com/owner/repo.git", "owner/repo", true},
		{"https://www.github.com/owner/repo/", "owner/repo", true},
		{"https://github.com/owner/repo?query=param", "owner/repo", true},
		{"https://github.com/owner/repo#fragment", "owner/repo", true},
		{"https://github.com/owner-with-dash/repo_with_underscore", "owner-with-dash/repo_with_underscore", true},
		{"https://gitlab.com/owner/repo", "", false},
		{"https://notgithub.com", "", false},
		{"invalid-url", "", false},
		{"https://github.com", "", false}, // Just domain, no owner/repo
		{"https://github.com/owner", "", false}, // Only owner, no repo
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			normalized, isGitHub := service.NormalizeGitHubURL(tc.input)
			if normalized != tc.expected {
				t.Errorf("Expected normalized URL '%s', got '%s'", tc.expected, normalized)
			}
			if isGitHub != tc.isGitHub {
				t.Errorf("Expected isGitHub %t, got %t", tc.isGitHub, isGitHub)
			}
		})
	}
}

func TestExtractService_Extract_FullPagination(t *testing.T) {
	// Create a mock source that returns all items in one call
	mockSource := &mockBookmarkSource{
		bookmarks: [][]domain.RawBookmark{
			{
				{ID: "1", Content: struct {
					URL         string `json:"url"`
					Title       string `json:"title"`
					Description string `json:"description"`
					HTMLContent string `json:"htmlContent"`
				}{URL: "https://github.com/a/b", Title: "Repo A", HTMLContent: ""}},
				{ID: "2", Content: struct {
					URL         string `json:"url"`
					Title       string `json:"title"`
					Description string `json:"description"`
					HTMLContent string `json:"htmlContent"`
				}{URL: "https://github.com/c/d", Title: "Repo C", HTMLContent: ""}},
				{ID: "3", Content: struct {
					URL         string `json:"url"`
					Title       string `json:"title"`
					Description string `json:"description"`
					HTMLContent string `json:"htmlContent"`
				}{URL: "https://github.com/e/f", Title: "Repo E", HTMLContent: ""}},
			},
		},
	}
	mockRepo := newMockRepoRepository()
	extractor := service.NewExtractor(mockSource, mockRepo)

	err := extractor.Extract(context.Background(), &mockReporter{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockRepo.repos) != 3 {
		t.Errorf("Expected 3 repos after pagination, got %d", len(mockRepo.repos))
	}
	if _, ok := mockRepo.repos["a/b"]; !ok {
		t.Errorf("Repo a/b not found")
	}
	if _, ok := mockRepo.repos["c/d"]; !ok {
		t.Errorf("Repo c/d not found")
	}
	if _, ok := mockRepo.repos["e/f"]; !ok {
		t.Errorf("Repo e/f not found")
	}
}