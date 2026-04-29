package logger

import "github.com/rs/zerolog/log"

// Fields is a convenience alias for structured logger fields.
type Fields map[string]any

func Info(msg string, fields Fields) {
	log.Info().Fields(map[string]any(fields)).Msg(msg)
}

func Warn(msg string, fields Fields) {
	log.Warn().Fields(map[string]any(fields)).Msg(msg)
}

func Error(msg string, fields Fields) {
	log.Error().Fields(map[string]any(fields)).Msg(msg)
}
