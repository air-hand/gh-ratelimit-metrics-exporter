package main

import (
	"github.com/google/go-github/v61/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

// https://docs.github.com/en/rest/rate-limit/rate-limit?apiVersion=2022-11-28#about-rate-limits
var (
	rateLimitCoreRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_rate_limit_core_remaining",
		Help: "The remaining number of requests to GitHub API",
	})
	rateLimitSearchRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_rate_limit_search_remaining",
		Help: "The remaining number of requests to GitHub API",
	})
	rateLimitCodeSearchRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_rate_limit_code_search_remaining",
		Help: "The remaining number of requests to GitHub API",
	})
	rateLimitGraphQLRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_rate_limit_graphql_remaining",
		Help: "The remaining number of requests to GitHub API",
	})
	rateLimitIntegrationManifestRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_rate_limit_integration_manifest_remaining",
		Help: "The remaining number of requests to GitHub API",
	})
	rateLimitDependencySnapshotsRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_rate_limit_dependency_snapshots_remaining",
		Help: "The remaining number of requests to GitHub API",
	})
	rateLimitCodeScanningUploadRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_rate_limit_code_scanning_upload_remaining",
		Help: "The remaining number of requests to GitHub API",
	})
	rateLimitActionsRunnerRegistrationRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_rate_limit_actions_runner_registration_remaining",
		Help: "The remaining number of requests to GitHub API",
	})
	rateLimitSourceImportRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_rate_limit_source_import_remaining",
		Help: "The remaining number of requests to GitHub API",
	})
	rateLimitSCIMRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_rate_limit_scim_remaining",
		Help: "The remaining number of requests to GitHub API",
	})
)

func initPrometheus() {
	prometheus.MustRegister(
		rateLimitCoreRemaining,
		rateLimitCodeSearchRemaining,
		rateLimitDependencySnapshotsRemaining,
		rateLimitActionsRunnerRegistrationRemaining,
		rateLimitCodeScanningUploadRemaining,
		rateLimitGraphQLRemaining,
		rateLimitIntegrationManifestRemaining,
		rateLimitSCIMRemaining,
		rateLimitSearchRemaining,
		rateLimitSourceImportRemaining,
	)
}

//go:generate moq -out gen_metrics_moq_test.go . RateLimitsFetcher

type RateLimitsFetcher interface {
	Fetch() (*github.RateLimits, error)
}

func fetchGitHubRateLimit(rlf RateLimitsFetcher, logger *zerolog.Logger) {
	rateLimits, err := rlf.Fetch()
	if err != nil {
		logger.Error().Err(err).Msg("Fail to fetch rate limits.")
		return
	}
	rateLimitCoreRemaining.Set(float64(rateLimits.Core.Remaining))
	rateLimitSearchRemaining.Set(float64(rateLimits.Search.Remaining))
	rateLimitCodeSearchRemaining.Set(float64(rateLimits.CodeSearch.Remaining))
	rateLimitGraphQLRemaining.Set(float64(rateLimits.GraphQL.Remaining))
	rateLimitIntegrationManifestRemaining.Set(float64(rateLimits.IntegrationManifest.Remaining))
	rateLimitDependencySnapshotsRemaining.Set(float64(rateLimits.DependencySnapshots.Remaining))
	rateLimitCodeScanningUploadRemaining.Set(float64(rateLimits.CodeScanningUpload.Remaining))
	rateLimitActionsRunnerRegistrationRemaining.Set(float64(rateLimits.ActionsRunnerRegistration.Remaining))
	rateLimitSourceImportRemaining.Set(float64(rateLimits.SourceImport.Remaining))
	rateLimitSCIMRemaining.Set(float64(rateLimits.SCIM.Remaining))
}
