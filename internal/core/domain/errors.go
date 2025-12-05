package domain

import "errors"

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrRepoNotFound      = errors.New("repository not found")
)
