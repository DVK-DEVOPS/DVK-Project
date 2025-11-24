package logging

import (
	"github.com/rs/zerolog"
	"os"
)

var Log = zerolog.New(os.Stdout).
	With().
	Timestamp().
	Logger()
