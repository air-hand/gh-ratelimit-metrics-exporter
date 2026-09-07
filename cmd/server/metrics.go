package main

import (
	"context"

	"github.com/google/go-github/v91/github"
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
	rateLimitAuditLogRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_rate_limit_audit_log_remaining",
		Help: "The remaining number of requests to GitHub API",
	})
	rateLimitDependencySBOMRemaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_rate_limit_dependency_sbom_remaining",
		Help: "The remaining number of requests to GitHub API",
	})
)

func initPrometheus() {
	prometheus.MustRegister(
		rateLimitAuditLogRemaining,
		rateLimitCoreRemaining,
		rateLimitCodeSearchRemaining,
		rateLimitDependencySnapshotsRemaining,
		rateLimitDependencySBOMRemaining,
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
	Fetch(ctx context.Context) (*github.RateLimits, error)
}

func fetchGitHubRateLimit(ctx context.Context, rlf RateLimitsFetcher, logger *zerolog.Logger) {
	rateLimits, err := rlf.Fetch(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Fail to fetch rate limits.")
		return
	}
	setRateLimitRemaining(rateLimitCoreRemaining, rateLimits.Core)
	setRateLimitRemaining(rateLimitSearchRemaining, rateLimits.Search)
	setRateLimitRemaining(rateLimitCodeSearchRemaining, rateLimits.CodeSearch)
	setRateLimitRemaining(rateLimitGraphQLRemaining, rateLimits.GraphQL)
	setRateLimitRemaining(rateLimitIntegrationManifestRemaining, rateLimits.IntegrationManifest)
	setRateLimitRemaining(rateLimitDependencySnapshotsRemaining, rateLimits.DependencySnapshots)
	setRateLimitRemaining(rateLimitCodeScanningUploadRemaining, rateLimits.CodeScanningUpload)
	setRateLimitRemaining(rateLimitActionsRunnerRegistrationRemaining, rateLimits.ActionsRunnerRegistration)
	setRateLimitRemaining(rateLimitSourceImportRemaining, rateLimits.SourceImport)
	setRateLimitRemaining(rateLimitSCIMRemaining, rateLimits.SCIM)
	setRateLimitRemaining(rateLimitAuditLogRemaining, rateLimits.AuditLog)
	setRateLimitRemaining(rateLimitDependencySBOMRemaining, rateLimits.DependencySBOM)
}

func setRateLimitRemaining(gauge prometheus.Gauge, rate *github.Rate) {
	if rate == nil {
		return
	}
	gauge.Set(float64(rate.Remaining))
}
