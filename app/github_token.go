package main

import (
	"cmp"
	"os"

	"github.com/google/go-github/v61/github"
	"github.com/rs/zerolog"
)

func newClientWithToken(logger *zerolog.Logger) *github.Client {
	token := cmp.Or(os.Getenv("GH_TOKEN"), os.Getenv("GITHUB_TOKEN"))
	if token != "" {
		logger.Info().Msg("Generate a new GitHub client with a token.")
		return github.NewClient(nil).WithAuthToken(token)
	}
	logger.Debug().Msg("Neighter GH_TOKEN nor GITHUB_TOKEN is set. Skip generating a new GitHub client with a token.")
	return nil
}
