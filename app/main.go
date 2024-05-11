package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.uber.org/dig"
)

func init() {
	initPrometheus()
}

func registerConstructors() *dig.Container {
	container := dig.New()
	if err := container.Provide(NewStdoutLogger); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(func() []newClientFunc {
		return newClientFuncs
	}); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(newClientWithEnv); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(newGitHubRateLimitsFetcher); err != nil {
		log.Fatal(err)
	}
	return container
}

func main() {
	c := registerConstructors()

	err := c.Invoke(func(rtf RateLimitsFetcher, logger *zerolog.Logger) {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		mux := http.NewServeMux()
		// TODO: impl health check endpoint.
		mux.Handle("/metrics", promhttp.Handler())
		// TODO: enable to change port.
		server := &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error().Err(err).Msg("Failed to start HTTP server.")
			}
		}()

		// fetch GitHub rate limit at the beginning
		fetchGitHubRateLimit(rtf, logger)

		go func() {
			// TODO: enable to change interval
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					logger.Info().Msg("Stop fetching GitHub rate limit.")
					return
				case <-ticker.C:
					fetchGitHubRateLimit(rtf, logger)
				}
			}
		}()

		<-ctx.Done()
		logger.Info().Msg("Received a signal to stop.")

		ctxS, cancelS := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelS()

		if err := server.Shutdown(ctxS); err != nil {
			logger.Fatal().Err(err).Msg("Failed to shutdown HTTP server.")
			return
		}
		logger.Info().Msg("Server shutdown")
	})

	if err != nil {
		log.Fatal(err)
	}
}
