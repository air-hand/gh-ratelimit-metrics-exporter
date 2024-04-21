package main

import (
	"testing"

	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
)

func TestFetchGitHubRateLimit_FailToFetch(t *testing.T) {
	rlf := &RateLimitsFetcherMock{
		FetchFunc: func() *github.RateLimits {
			return nil
		},
	}

	logger := NewNullLogger()

	assert.NotPanics(t, func() {
		fetchGitHubRateLimit(rlf, logger)
	})
}

// TODO: write test for fetchGitHubRateLimit
