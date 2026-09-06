package main

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/google/go-github/v91/github"
	"github.com/jferrl/go-githubauth"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

// newClientOptsWithGitHubApp creates a new GitHub client with GitHub App.
func newClientOptsWithGitHubApp(logger *zerolog.Logger) ([]github.ClientOptionsFunc, error) {
	ghAppClientID := os.Getenv("GH_APP_CLIENT_ID")
	ghInstallationID := os.Getenv("GH_INSTALLATION_ID")
	ghPrivateKey := os.Getenv("GH_PRIVATE_KEY")

	if ghAppClientID == "" && ghInstallationID == "" && ghPrivateKey == "" {
		logger.Debug().Msg("GitHub App envs are not set. Skip auth as app.")
		return nil, nil
	}

	if ghAppClientID == "" || ghInstallationID == "" || ghPrivateKey == "" {
		return nil, errors.New("GitHub App auth is partially configured")
	}

	ghInstallationIDInt64, err := strconv.ParseInt(ghInstallationID, 10, 64)
	if err != nil {
		logger.Error().Msg("Failed to parse GH_INSTALLATION_ID as integer.")
		return nil, err
	}

	client, err := newGitHubAppHTTPClient(ghAppClientID, ghInstallationIDInt64, []byte(ghPrivateKey))
	if err != nil {
		logger.Error().Msg("Failed to new HTTP client with GitHub App.")
		return nil, err
	}
	logger.Info().Msg("Generate a new GitHub client with GitHub App.")
	return []github.ClientOptionsFunc{
		github.WithHTTPClient(client),
	}, nil
}

func newGitHubAppHTTPClient(appClientID string, installationID int64, privateKey []byte) (*http.Client, error) {
	appTokenSource, err := githubauth.NewApplicationTokenSource(appClientID, privateKey)
	if err != nil {
		return nil, err
	}

	installationTokenSource := githubauth.NewInstallationTokenSource(
		installationID,
		appTokenSource,
		githubauth.WithInstallationTokenOptions(emptyInstallationTokenOptions()),
	)

	return oauth2.NewClient(context.Background(), installationTokenSource), nil
}

func emptyInstallationTokenOptions() *githubauth.InstallationTokenOptions {
	return &githubauth.InstallationTokenOptions{
		Permissions: &githubauth.InstallationPermissions{},
	}
}
