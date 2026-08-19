package pkg

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger membungkus *zap.Logger supaya bisa dipasang lewat pkg.Package dan
// dipakai lintas service. Method zap (Info/Warn/Error/Debug/With/...) tersedia
// langsung karena embedding.
type Logger struct {
	*zap.Logger
}

func NewLogger() (*Logger, error) {
	var cfg zap.Config
	if strings.EqualFold(os.Getenv("APP_ENV"), "production") {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Jangan lampirkan stacktrace di log error biasa; cukup pesan + field-nya.
	cfg.DisableStacktrace = true

	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		if parsed, err := zapcore.ParseLevel(lvl); err == nil {
			cfg.Level = zap.NewAtomicLevelAt(parsed)
		}
	}

	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: l}, nil
}
