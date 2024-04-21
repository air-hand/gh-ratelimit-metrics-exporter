package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.uber.org/dig"
)

func init() {
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

func registerConstructors() *dig.Container {
	container := dig.New()
	container.Provide(NewStdoutLogger)
	container.Provide(func() []newClientFunc {
		return newClientFuncs
	})
	container.Provide(newClientWithEnv)
	container.Provide(newGitHubRateLimitsFetcher)
	return container
}

func main() {
	c := registerConstructors()

	err := c.Invoke(func(rtf RateLimitsFetcher, logger *zerolog.Logger) {
		go func() {
			for {
				fetchGitHubRateLimit(rtf, logger)
				// TODO: enable to change interval
				time.Sleep(5 * time.Minute)
			}
		}()
	})

	if err != nil {
		panic(err)
	}

	http.Handle("/metrics", promhttp.Handler())
	// TODO: impl health check endpoint.
	// TODO: enable to change port.
	http.ListenAndServe(":8080", nil)
}
