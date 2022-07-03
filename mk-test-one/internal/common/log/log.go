package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aliakbariaa1996/mk-test-one/internal/common/configx"
)

type Option func(*Logger)

type Logger struct {
	zapLogger    *zap.SugaredLogger
	sentryOption sentryOption
}

type sentryOption struct {
	sentryDsn    string
	sentryTags   map[string]string
	sentryFields []zapcore.Field
}

func WithSentry(dsn string, tags map[string]string, fields ...zapcore.Field) Option {
	return func(logger *Logger) {
		logger.sentryOption.sentryDsn = dsn
		logger.sentryOption.sentryTags = tags
		logger.sentryOption.sentryFields = fields
	}
}

// NewTestLogger return instance of Logger that discards all output.
func NewTestLogger() *Logger {
	return &Logger{
		zapLogger: zap.NewNop().Sugar(),
	}
}

func New(mode, serviceName string, opts ...Option) (*Logger, error) {
	zapLogger, err := newZap(mode)
	if err != nil {
		return nil, fmt.Errorf("failed to create zap sugared: %w", err)
	}

	logger := &Logger{
		zapLogger: zapLogger.Sugar(),
	}
	for _, opt := range opts {
		opt(logger)
	}

	if logger.sentryOption.sentryDsn == "" || mode != configx.ModeProd {
		return logger, nil
	}

	sentryCore, err := newSentryCore(newSentryOptions(logger.sentryOption.sentryDsn, mode, serviceName))
	if err != nil {
		return nil, fmt.Errorf("failed to init sentry core: %w", err)
	}

	core := zapcore.NewTee(zapLogger.Core(), sentryCore)
	logger.zapLogger = zap.New(core).Sugar()

	return logger, nil
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.zapLogger.Infow(msg, fieldsToInterface(fields)...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.zapLogger.Errorw(msg, fieldsToInterface(fields)...)
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.zapLogger.Debugw(msg, fieldsToInterface(fields)...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.zapLogger.Warnw(msg, fieldsToInterface(fields)...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.zapLogger.Infof(template, args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.zapLogger.Errorf(template, args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.zapLogger.Debugf(template, args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.zapLogger.Warnf(template, args...)
}

func newZap(mode string) (*zap.Logger, error) {
	opts := []zap.Option{zap.AddCallerSkip(1)}

	if mode == modeDev {
		return zap.NewDevelopment(opts...)
	}

	return zap.NewProduction(opts...)
}
