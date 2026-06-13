package slogcolored

import (
	"log/slog"
	"os"

	"go-vents/internal/domain/types"

	"github.com/lmittmann/tint"
)

var groupPrefix = ""

type ColoredSLogger struct {
	logger *slog.Logger
}

func NewColoredSLogger() types.Logger {

	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    false,
		AddSource:  false,
	})
	logger := slog.New(handler)

	return &ColoredSLogger{
		logger: logger,
	}
}

func (s *ColoredSLogger) OpenGroup(name string) {
	s.logger.Info("----------------------")
	s.logger.Info("📂 - start of " + name)
	s.logger.Info("----------------------")

	groupPrefix = name + " > "
}

func (s *ColoredSLogger) CloseGroup(name string) {
	s.logger.Info("----------------------")
	s.logger.Info("📁 - end of " + name)
	s.logger.Info("----------------------")

	groupPrefix = ""
}

func (s *ColoredSLogger) Info(message string) {
	s.logger.Info("ℹ️  ➜  " + groupPrefix + message)
}

func (s *ColoredSLogger) Error(message string) {
	s.logger.Warn("❌ ➜  " + groupPrefix + message)
}

func (s *ColoredSLogger) Fatal(message string) {
	s.logger.Error("💀 ➜  " + groupPrefix + message)
	os.Exit(1)
}
