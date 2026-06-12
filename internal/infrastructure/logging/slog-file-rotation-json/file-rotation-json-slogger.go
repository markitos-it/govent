package slogfilerotationjson

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"govent/internal/domain/types"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"gopkg.in/natefinch/lumberjack.v2"
)

type FileRotationJSONSLogger struct {
	logger *slog.Logger
}

func NewFileRotationJSONSLogger() types.Logger {
	logDir := "logs"
	_ = os.MkdirAll(logDir, 0750)

	fileRotator := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "production.log"),
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	logDestination := io.MultiWriter(os.Stdout, fileRotator)

	baseHandler := slog.NewJSONHandler(logDestination, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(baseHandler)

	return &FileRotationJSONSLogger{
		logger: logger,
	}
}

func (s *FileRotationJSONSLogger) GetInternalLogger() *slog.Logger {
	return s.logger
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (s *FileRotationJSONSLogger) OpenGroup(name string) {
	s.logger.Info("----- start of group")
}

func (s *FileRotationJSONSLogger) CloseGroup(name string) {
	s.logger.Info("----- end of group")
}

func (s *FileRotationJSONSLogger) Info(message string) {
	s.logger.Info(message)
}

func (s *FileRotationJSONSLogger) Error(message string) {
	s.logger.Error(message)
}

func (s *FileRotationJSONSLogger) Fatal(message string) {
	s.logger.Error(message, "fatal", true)
	os.Exit(1)
}
