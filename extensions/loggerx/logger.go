package loggerx

import (
	"os"
	"runtime"
	"time"

	"github.com/rs/zerolog"
)

// NewDevelopment creates a new zerolog Logger configured for development environment.
// The logger outputs to the console with human-readable formatting and includes
// additional context such as the Go runtime version, process ID, and caller information.
//
// Returns:
//   - zerolog.Logger: A configured zerolog Logger instance for development purposes.
func NewDevelopment() zerolog.Logger {
	goVersion := runtime.Version()
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", goVersion).
		Logger()
}
