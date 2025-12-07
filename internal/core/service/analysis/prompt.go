package analysis

import (
	"encoding/json"
	"fmt"

	"github.com/brianluby/karakeep-extractor/internal/core/domain"
)

// BuildMessages constructs the chat messages for the LLM analysis.
func BuildMessages(query string, repos []domain.ExtractedRepo) ([]domain.Message, error) {
	contextJSON, err := serializeRepos(repos)
	if err != nil {
		return nil, err
	}

	systemContent := fmt.Sprintf(`You are an expert software engineering assistant.
You are analyzing a curated list of GitHub repositories provided in JSON format.
Your goal is to answer the user's question based on the provided repository data.

Repository Data:
%s

Instructions:
- Base your answer strictly on the provided data.
- If the answer cannot be found in the data, say so.
- Be concise and helpful.
`, contextJSON)

	return []domain.Message{
		{Role: "system", Content: systemContent},
		{Role: "user", Content: query},
	}, nil
}

func serializeRepos(repos []domain.ExtractedRepo) (string, error) {
	var contexts []domain.RepositoryContext
	for _, r := range repos {
		ctx := domain.RepositoryContext{
			Name: r.RepoID,
			URL:  r.URL,
		}
		if r.Description != nil {
			ctx.Description = *r.Description
		}
		if r.Language != nil {
			ctx.Language = *r.Language
		}
		if r.Stars != nil {
			ctx.Stars = *r.Stars
		}
		if r.Forks != nil {
			ctx.Forks = *r.Forks
		}
		if r.LastPushedAt != nil {
			ctx.LastUpdated = r.LastPushedAt.Format("2006-01-02")
		}
		// Tags not available in ExtractedRepo yet
		
		contexts = append(contexts, ctx)
	}

	// Use Marshal without indent to save tokens? Or Indent for readability?
	// Indent uses more tokens. We should probably use compact.
	data, err := json.Marshal(contexts)
	if err != nil {
		return "", fmt.Errorf("failed to serialize repos: %w", err)
	}
	return string(data), nil
}
