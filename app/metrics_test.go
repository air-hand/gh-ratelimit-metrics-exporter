package main

import (
	"fmt"
	"testing"

	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
)

func TestFetchGitHubRateLimit_FailToFetch(t *testing.T) {
	rlf := &RateLimitsFetcherMock{
		FetchFunc: func() (*github.RateLimits, error) {
			return nil, fmt.Errorf("some errors")
		},
	}

	logger := NewNullLogger()

	assert.NotPanics(t, func() {
		fetchGitHubRateLimit(rlf, logger)
	})
}

// TODO: write test for fetchGitHubRateLimit
