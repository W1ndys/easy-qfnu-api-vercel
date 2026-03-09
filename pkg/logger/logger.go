package logger

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level      string
	Format     string
	File       string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

var (
	mu     sync.RWMutex
	global = zap.NewNop()
	sugar  = global.Sugar()
)

func LoadConfigFromEnv() Config {
	return Config{
		Level:      normalizeString(os.Getenv("LOG_LEVEL"), "info"),
		Format:     normalizeString(os.Getenv("LOG_FORMAT"), "console"),
		File:       strings.TrimSpace(os.Getenv("LOG_FILE")),
		MaxSize:    parseIntEnv("LOG_MAX_SIZE", 100),
		MaxBackups: parseIntEnv("LOG_MAX_BACKUPS", 30),
		MaxAge:     parseIntEnv("LOG_MAX_AGE", 30),
	}
}

func Init() error {
	cfg := LoadConfigFromEnv()
	logger, err := newLogger(cfg)
	if err != nil {
		return err
	}

	mu.Lock()
	global = logger
	sugar = logger.Sugar()
	mu.Unlock()

	zap.ReplaceGlobals(logger)
	slog.SetDefault(slog.New(&zapSlogHandler{logger: logger}))

	return nil
}

func L() *zap.Logger {
	mu.RLock()
	defer mu.RUnlock()
	return global
}

func S() *zap.SugaredLogger {
	mu.RLock()
	defer mu.RUnlock()
	return sugar
}

func Sync() {
	_ = L().Sync()
}

func newLogger(cfg Config) (*zap.Logger, error) {
	level := parseLevel(cfg.Level)
	encoder := newEncoder(cfg.Format)
	writer := newWriteSyncer(cfg)

	core := zapcore.NewCore(encoder, writer, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}

func newEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"

	if strings.EqualFold(format, "json") {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func newWriteSyncer(cfg Config) zapcore.WriteSyncer {
	writers := []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}
	if cfg.File != "" {
		writers = append(writers, zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.File,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   false,
		}))
	}

	return zapcore.NewMultiWriteSyncer(writers...)
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return zapcore.DebugLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func normalizeString(value string, defaultValue string) string {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	if trimmed == "" {
		return defaultValue
	}
	return trimmed
}

func parseIntEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

type zapSlogHandler struct {
	logger *zap.Logger
	attrs  []slog.Attr
	groups []string
}

func (h *zapSlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.logger.Core().Enabled(slogLevelToZap(level))
}

func (h *zapSlogHandler) Handle(_ context.Context, record slog.Record) error {
	fields := make([]zap.Field, 0, len(h.attrs)+record.NumAttrs())
	fields = append(fields, attrsToFields(h.groups, h.attrs)...)
	record.Attrs(func(attr slog.Attr) bool {
		fields = append(fields, attrsToFields(h.groups, []slog.Attr{attr})...)
		return true
	})

	switch slogLevelToZap(record.Level) {
	case zapcore.DebugLevel:
		h.logger.Debug(record.Message, fields...)
	case zapcore.WarnLevel:
		h.logger.Warn(record.Message, fields...)
	case zapcore.ErrorLevel:
		h.logger.Error(record.Message, fields...)
	default:
		h.logger.Info(record.Message, fields...)
	}

	return nil
}

func (h *zapSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	cloned := &zapSlogHandler{
		logger: h.logger,
		attrs:  append(append([]slog.Attr(nil), h.attrs...), attrs...),
		groups: append([]string(nil), h.groups...),
	}
	return cloned
}

func (h *zapSlogHandler) WithGroup(name string) slog.Handler {
	cloned := &zapSlogHandler{
		logger: h.logger,
		attrs:  append([]slog.Attr(nil), h.attrs...),
		groups: append(append([]string(nil), h.groups...), name),
	}
	return cloned
}

func attrsToFields(groups []string, attrs []slog.Attr) []zap.Field {
	fields := make([]zap.Field, 0, len(attrs))
	for _, attr := range attrs {
		if attr.Equal(slog.Attr{}) {
			continue
		}
		fields = append(fields, attrToField(groups, attr)...)
	}
	return fields
}

func attrToField(groups []string, attr slog.Attr) []zap.Field {
	attr.Value = attr.Value.Resolve()

	if attr.Value.Kind() == slog.KindGroup {
		childGroups := groups
		if attr.Key != "" {
			childGroups = append(append([]string(nil), groups...), attr.Key)
		}
		return attrsToFields(childGroups, attr.Value.Group())
	}

	key := attr.Key
	if key == "" {
		key = "value"
	}
	if len(groups) > 0 {
		key = strings.Join(append(append([]string(nil), groups...), key), ".")
	}

	return []zap.Field{zap.Any(key, slogValueToAny(attr.Value))}
}

func slogLevelToZap(level slog.Level) zapcore.Level {
	switch {
	case level <= slog.LevelDebug:
		return zapcore.DebugLevel
	case level < slog.LevelWarn:
		return zapcore.InfoLevel
	case level < slog.LevelError:
		return zapcore.WarnLevel
	default:
		return zapcore.ErrorLevel
	}
}

func slogValueToAny(value slog.Value) any {
	switch value.Kind() {
	case slog.KindBool:
		return value.Bool()
	case slog.KindDuration:
		return value.Duration()
	case slog.KindFloat64:
		return value.Float64()
	case slog.KindInt64:
		return value.Int64()
	case slog.KindString:
		return value.String()
	case slog.KindTime:
		return value.Time()
	case slog.KindUint64:
		return value.Uint64()
	case slog.KindAny:
		return value.Any()
	default:
		return value.Any()
	}
}
