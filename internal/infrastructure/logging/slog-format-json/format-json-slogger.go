package slogformatjson

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"govent/internal/domain/types"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

type PrettyJSONHandler struct {
	slog.Handler
}

func (h *PrettyJSONHandler) Handle(ctx context.Context, r slog.Record) error {

	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	var buf bytes.Buffer

	err := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: r.Level,
	}).Handle(ctx, r)

	if err != nil {
		return err
	}

	var indentedBuf bytes.Buffer
	err = json.Indent(&indentedBuf, buf.Bytes(), "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(indentedBuf.String())
	return nil
}

type JSONSLogger struct {
	logger *slog.Logger
}

func NewJSONSLogger() types.Logger {
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(&PrettyJSONHandler{Handler: baseHandler})

	return &JSONSLogger{
		logger: logger,
	}
}

func (s *JSONSLogger) GetInternalLogger() *slog.Logger {
	return s.logger
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (s *JSONSLogger) OpenGroup(name string) {
	s.logger.Info("start of group", "group_name", name)
}

func (s *JSONSLogger) CloseGroup(name string) {
	s.logger.Info("end of group", "group_name", name)
}

func (s *JSONSLogger) Info(message string) {
	s.logger.Info(message)
}

func (s *JSONSLogger) Error(message string) {
	s.logger.Warn(message)
}

func (s *JSONSLogger) Fatal(message string) {
	s.logger.Error(message)
	os.Exit(1)
}
