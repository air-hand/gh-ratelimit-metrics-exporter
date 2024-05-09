package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v61/github"
	"github.com/rs/zerolog"
)

// newClientWithGitHubApp creates a new GitHub client with GitHub App.
func newClientWithGitHubApp(logger *zerolog.Logger) *github.Client {
	gh_app_id := os.Getenv("GH_APP_ID")
	gh_installation_id := os.Getenv("GH_INSTALLATION_ID")
	gh_private_key := os.Getenv("GH_PRIVATE_KEY")

	if gh_app_id == "" || gh_installation_id == "" || gh_private_key == "" {
		logger.Debug().Msg("Neighter GH_APP_ID nor GH_INSTALLATION_ID nor GH_PRIVATE_KEY is set. Skip generating a new GitHub client with GitHub App.")
		return nil
	}

	gh_app_id_int64, err := strconv.ParseInt(gh_app_id, 10, 64)
	if err != nil {
		logger.Error().Msg("Failed to parse GH_APP_ID as integer.")
		return nil
	}
	gh_installation_id_int64, err := strconv.ParseInt(gh_installation_id, 10, 64)
	if err != nil {
		logger.Error().Msg("Failed to parse GH_INSTALLATION_ID as integer.")
		return nil
	}
	tr, err := ghinstallation.New(http.DefaultTransport, gh_app_id_int64, gh_installation_id_int64, []byte(gh_private_key))
	if err != nil {
		logger.Error().Msg("Failed to new transport with GitHub App.")
		return nil
	}
	logger.Info().Msg("Generate a new GitHub client with GitHub App.")
	return github.NewClient(&http.Client{Transport: tr})
}
