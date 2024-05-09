package main

import (
	"context"
	"log"
	"net/http"
	"os"
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
	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	err := c.Invoke(func(rtf RateLimitsFetcher, logger *zerolog.Logger) {
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

		// ここはメインスレッドで実行しないとだめかもしれない メインスレッドが終了しない...
		go func() {
			<-quit
			logger.Info().Msg("Received a signal to stop.")
			cancel()

			ctxS, cancelS := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelS()

			if err := server.Shutdown(ctxS); err != nil {
				logger.Error().Err(err).Msg("Failed to shutdown HTTP server.")
				return
			}
			logger.Info().Msg("Server shutdown")
		}()
	})

	if err != nil {
		log.Panic(err)
	}

	// blocking
	select {}
}
