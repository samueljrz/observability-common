package log

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/garden/observability-commons/config"
	"github.com/garden/observability-commons/util"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	instrumentationName = "github.com/garden/observability-commons"
)

type OTLPLogger struct {
	logger *zap.Logger
	cfg    config.Config
	tracer trace.Tracer
}

func NewOTLPLogger(cfg config.Config) (*OTLPLogger, error) {

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stack_trace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var level zap.AtomicLevel
	if cfg.Mode != config.Production {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	var core zapcore.Core
	switch cfg.Mode {
	case config.Noop:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.NewFile(0, os.DevNull)),
			level,
		)
	case config.Local:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
	case config.Debug, config.Development, config.Production:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
	default:
		return nil, fmt.Errorf("unknown mode: %v", cfg.Mode)
	}

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	tracer := trace.NewNoopTracerProvider().Tracer(instrumentationName)

	return &OTLPLogger{
		logger: logger,
		cfg:    cfg,
		tracer: tracer,
	}, nil
}

func (log *OTLPLogger) Debug(logEntry *Entry) {
	log.logWithLevel(logEntry, zap.DebugLevel)
}

func (log *OTLPLogger) Info(logEntry *Entry) {
	log.logWithLevel(logEntry, zap.InfoLevel)
}

func (log *OTLPLogger) Warn(logEntry *Entry) {
	logEntry.stacktrace = string(debug.Stack())
	log.logWithLevel(logEntry, zap.WarnLevel)
}

func (log *OTLPLogger) Error(logEntry *Entry) {
	logEntry.stacktrace = string(debug.Stack())
	log.logWithLevel(logEntry, zap.ErrorLevel)
}

func (log *OTLPLogger) Fatal(logEntry *Entry) {
	logEntry.stacktrace = string(debug.Stack())
	log.logWithLevel(logEntry, zap.FatalLevel)
}

func (log *OTLPLogger) Close() error {
	return log.logger.Sync()
}

func (log *OTLPLogger) logWithLevel(logEntry *Entry, level zapcore.Level) {
	fields := log.generateOTLPFields(logEntry)

	go func() {
		switch level {
		case zap.DebugLevel:
			log.logger.Debug(logEntry.Message, fields...)
		case zap.InfoLevel:
			log.logger.Info(logEntry.Message, fields...)
		case zap.WarnLevel:
			log.logger.Warn(logEntry.Message, fields...)
		case zap.ErrorLevel:
			log.logger.Error(logEntry.Message, fields...)
		case zap.FatalLevel:
			log.logger.Fatal(logEntry.Message, fields...)
		}
	}()
}

func (log *OTLPLogger) generateOTLPFields(logEntry *Entry) []zap.Field {
	fields := []zap.Field{
		zap.String("service.name", log.cfg.Service.Name),
		zap.String("service.version", log.cfg.Service.Version),
		zap.String("host.name", log.cfg.GetHostname()),
		zap.String("component", logEntry.Component),
		zap.String("operation", logEntry.Operation),
		zap.Time("timestamp", time.Now()),
	}

	if logEntry.Err != nil {
		fields = append(fields, zap.Error(logEntry.Err))
	}

	if logEntry.stacktrace != "" {
		stacktraceHash := util.MD5Hash([]byte(logEntry.stacktrace))
		fields = append(fields, zap.String("stacktrace.hash", stacktraceHash))
		fields = append(fields, zap.String("stacktrace", logEntry.stacktrace))
	}

	for key, value := range logEntry.Fields {
		fields = append(fields, zap.String(key, value))
	}

	if log.cfg.DefaultFields != nil {
		for key, value := range *log.cfg.DefaultFields {
			fields = append(fields, zap.String(key, value))
		}
	}

	return fields
}
