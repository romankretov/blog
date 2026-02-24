package observability

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(env string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	level := zerolog.InfoLevel
	if env == "dev" {
		level = zerolog.DebugLevel
	}
	return zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()
}
