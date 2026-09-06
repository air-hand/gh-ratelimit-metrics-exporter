package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-github/v91/github"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFetchGitHubRateLimit_FailToFetch(t *testing.T) {
	rlf := &RateLimitsFetcherMock{
		FetchFunc: func(_ context.Context) (*github.RateLimits, error) {
			return nil, fmt.Errorf("some errors")
		},
	}

	logger := NewNullLogger()
	ctx := context.Background()

	assert.NotPanics(t, func() {
		fetchGitHubRateLimit(ctx, rlf, logger)
	})
}

func TestFetchGitHubRateLimit_SetRemainingMetrics(t *testing.T) {
	rlf := &RateLimitsFetcherMock{
		FetchFunc: func(_ context.Context) (*github.RateLimits, error) {
			return &github.RateLimits{
				Core:                      &github.Rate{Remaining: 1},
				Search:                    &github.Rate{Remaining: 2},
				CodeSearch:                &github.Rate{Remaining: 3},
				GraphQL:                   &github.Rate{Remaining: 4},
				IntegrationManifest:       &github.Rate{Remaining: 5},
				DependencySnapshots:       &github.Rate{Remaining: 6},
				CodeScanningUpload:        &github.Rate{Remaining: 7},
				ActionsRunnerRegistration: &github.Rate{Remaining: 8},
				SourceImport:              &github.Rate{Remaining: 9},
				SCIM:                      &github.Rate{Remaining: 10},
				AuditLog:                  &github.Rate{Remaining: 11},
				DependencySBOM:            &github.Rate{Remaining: 12},
			}, nil
		},
	}

	logger := NewNullLogger()
	ctx := context.Background()

	fetchGitHubRateLimit(ctx, rlf, logger)

	assert.Equal(t, float64(1), testutil.ToFloat64(rateLimitCoreRemaining))
	assert.Equal(t, float64(2), testutil.ToFloat64(rateLimitSearchRemaining))
	assert.Equal(t, float64(3), testutil.ToFloat64(rateLimitCodeSearchRemaining))
	assert.Equal(t, float64(4), testutil.ToFloat64(rateLimitGraphQLRemaining))
	assert.Equal(t, float64(5), testutil.ToFloat64(rateLimitIntegrationManifestRemaining))
	assert.Equal(t, float64(6), testutil.ToFloat64(rateLimitDependencySnapshotsRemaining))
	assert.Equal(t, float64(7), testutil.ToFloat64(rateLimitCodeScanningUploadRemaining))
	assert.Equal(t, float64(8), testutil.ToFloat64(rateLimitActionsRunnerRegistrationRemaining))
	assert.Equal(t, float64(9), testutil.ToFloat64(rateLimitSourceImportRemaining))
	assert.Equal(t, float64(10), testutil.ToFloat64(rateLimitSCIMRemaining))
	assert.Equal(t, float64(11), testutil.ToFloat64(rateLimitAuditLogRemaining))
	assert.Equal(t, float64(12), testutil.ToFloat64(rateLimitDependencySBOMRemaining))
}

func TestFetchGitHubRateLimit_SkipNilRateLimit(t *testing.T) {
	rlf := &RateLimitsFetcherMock{
		FetchFunc: func(_ context.Context) (*github.RateLimits, error) {
			return &github.RateLimits{}, nil
		},
	}

	logger := NewNullLogger()
	ctx := context.Background()

	assert.NotPanics(t, func() {
		fetchGitHubRateLimit(ctx, rlf, logger)
	})
}
