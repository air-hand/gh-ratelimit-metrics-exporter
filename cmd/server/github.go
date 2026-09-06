package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/go-github/v91/github"
	"github.com/rs/zerolog"
)

type newClientOptsFunc func(*zerolog.Logger) ([]github.ClientOptionsFunc, error)

var newClientFuncs = []newClientOptsFunc{
	newClientOptsWithGitHubApp,
	newClientOptsWithToken,
}

func newClientFromEnv(funcs []newClientOptsFunc, logger *zerolog.Logger) (*github.Client, error) {
	for _, f := range funcs {
		opts, err := f(logger)
		if err != nil {
			return nil, errors.Wrap(err, "failed to newClientFromEnv")
		}
		if len(opts) > 0 {
			opts = append(opts, github.WithTimeout(15*time.Second))
			return github.NewClient(opts...)
		}
	}
	return nil, errors.New("no GitHub authentication is configured")
}

type GitHubRateLimitsFetcher struct {
	client *github.Client
	logger *zerolog.Logger
}

// explicit compile error check
var _ RateLimitsFetcher = (*GitHubRateLimitsFetcher)(nil)

func (g *GitHubRateLimitsFetcher) Fetch(ctx context.Context) (*github.RateLimits, error) {
	rateLimits, res, err := g.client.RateLimit.Get(ctx)
	if err != nil {
		g.logger.Error().Err(err).Msgf("Err: %v", res)
		return nil, fmt.Errorf("%w", err)
	}
	return rateLimits, nil
}

func newGitHubRateLimitsFetcher(client *github.Client, logger *zerolog.Logger) RateLimitsFetcher {
	return &GitHubRateLimitsFetcher{client: client, logger: logger}
}
