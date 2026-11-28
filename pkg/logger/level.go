package logger

import "github.com/rs/zerolog"

func SetLogLevel(logLevel string) {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel // fallback
	}
	zerolog.SetGlobalLevel(level)
}
