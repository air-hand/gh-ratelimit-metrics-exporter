package main

import (
	"cmp"
	"os"

	"github.com/google/go-github/v91/github"
	"github.com/rs/zerolog"
)

func newClientOptsWithToken(logger *zerolog.Logger) ([]github.ClientOptionsFunc, error) {
	token := cmp.Or(os.Getenv("GH_TOKEN"), os.Getenv("GITHUB_TOKEN"))
	if token == "" {
		logger.Debug().Msg("Neither GH_TOKEN nor GITHUB_TOKEN is set. Skip generating a new GitHub client with a token.")
		return nil, nil
	}
	logger.Info().Msg("Generate a new GitHub client with a token.")
	return []github.ClientOptionsFunc{
		github.WithAuthToken(token),
	}, nil
}
