package main

import (
	"os"

	"github.com/google/go-github/v61/github"
	"github.com/rs/zerolog"
)

func newClientWithToken(logger *zerolog.Logger) *github.Client {
	if token := os.Getenv("GH_TOKEN"); token != "" {
		logger.Info().Msg("Generate a new GitHub client with GH_TOKEN.")
		return github.NewClient(nil).WithAuthToken(token)
	}
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		logger.Info().Msg("Generate a new GitHub client with GITHUB_TOKEN.")
		return github.NewClient(nil).WithAuthToken(token)
	}
	logger.Debug().Msg("GH_TOKEN or GITHUB_TOKEN is not set. Skip generating a new GitHub client with token.")
	return nil
}
