package httpx

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

// Run starts the HTTP server and listens for shutdown signals.
// It gracefully shuts down the server when an interrupt or terminate signal is received.
//
// Parameters:
//   - ctx: The context for managing the lifecycle of the server and logging.
//   - server: The HTTP server instance to be started and managed.
func Run(ctx context.Context, server *http.Server) {
	stopped := make(chan struct{})
	go func() {
		log.Ctx(ctx).Info().Str("address", server.Addr).
			Float64("read_timeout_sec", server.ReadTimeout.Seconds()).
			Float64("write_timeout_sec", server.WriteTimeout.Seconds()).
			Msg("server started")

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Ctx(ctx).Error().Err(err).Msg("http server failed")
			close(stopped)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-signals:
		log.Ctx(ctx).Info().Str("signal", sig.String()).Msg("captured signal")
	case <-stopped:
		log.Ctx(ctx).Info().Msg("server encountered an error and is shutting down")
	}

	signal.Stop(signals)

	log.Ctx(ctx).Info().Msg("shutting down HTTP server")
	c, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*5))
	defer cancel()

	if err := server.Shutdown(c); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error shutting down HTTP server")
	}
}
