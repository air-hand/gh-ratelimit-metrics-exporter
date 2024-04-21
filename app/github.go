package main

import (
	"context"

	"github.com/google/go-github/v61/github"
	"github.com/rs/zerolog"
)

type newClientFunc func(*zerolog.Logger) *github.Client

var newClientFuncs = []newClientFunc{
	newClientWithGitHubApp,
	newClientWithToken,
}

func newClientWithEnv(funcs []newClientFunc, logger *zerolog.Logger) *github.Client {
	for _, f := range funcs {
		if gh := f(logger); gh != nil {
			return gh
		}
	}
	// fallback
	logger.Info().Msg("fallback to No Auth client.")
	return github.NewClient(nil)
}

type GitHubRateLimitsFetcher struct {
	client *github.Client
	logger *zerolog.Logger
}

// explicit compile error check
var _ RateLimitsFetcher = (*GitHubRateLimitsFetcher)(nil)

func (g *GitHubRateLimitsFetcher) Fetch() *github.RateLimits {
	rate_limits, res, err := g.client.RateLimit.Get(context.Background())
	if err != nil {
		g.logger.Error().Msgf("Err: %v %v", res, err)
		return nil
	}
	return rate_limits
}

func newGitHubRateLimitsFetcher(client *github.Client, logger *zerolog.Logger) RateLimitsFetcher {
	return &GitHubRateLimitsFetcher{client: client, logger: logger}
}
