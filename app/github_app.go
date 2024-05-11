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
	ghAppID := os.Getenv("GH_APP_ID")
	ghInstallationID := os.Getenv("GH_INSTALLATION_ID")
	ghPrivateKey := os.Getenv("GH_PRIVATE_KEY")

	if ghAppID == "" || ghInstallationID == "" || ghPrivateKey == "" {
		logger.Debug().Msg("Neighter GH_APP_ID nor GH_INSTALLATION_ID nor GH_PRIVATE_KEY is set. Skip generating a new GitHub client with GitHub App.")
		return nil
	}

	ghAppIDInt64, err := strconv.ParseInt(ghAppID, 10, 64)
	if err != nil {
		logger.Error().Msg("Failed to parse GH_APP_ID as integer.")
		return nil
	}
	ghInstallationIDInt64, err := strconv.ParseInt(ghInstallationID, 10, 64)
	if err != nil {
		logger.Error().Msg("Failed to parse GH_INSTALLATION_ID as integer.")
		return nil
	}
	tr, err := ghinstallation.New(http.DefaultTransport, ghAppIDInt64, ghInstallationIDInt64, []byte(ghPrivateKey))
	if err != nil {
		logger.Error().Msg("Failed to new transport with GitHub App.")
		return nil
	}
	logger.Info().Msg("Generate a new GitHub client with GitHub App.")
	return github.NewClient(&http.Client{Transport: tr})
}
