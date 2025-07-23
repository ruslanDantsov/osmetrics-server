package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

// Log — глобальная переменная логгера, доступная для использования в приложении.
// По умолчанию инициализируется как zap.NewNop() (ничего не логирует).
// После вызова функции Initialized() содержит настроенный логгер.
var Log = zap.NewNop()

// Initialized инициализирует глобальный логгер Log на основе заданного уровня логирования.
// Уровень должен быть строкой: "debug", "info", "warn", "error" и т.д.
// Также применяется формат времени и имя сервиса по умолчанию.
// Возвращает ошибку в случае некорректного уровня или ошибки конфигурации.
func Initialized(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	withCustomTimeLayout("2006-01-02 15:04:05")(&cfg)
	withServiceName("Server app")(&cfg)

	configuredLogger, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = configuredLogger
	return nil
}

type LoggerOption func(*zap.Config)

func withCustomTimeLayout(layout string) LoggerOption {
	return func(cfg *zap.Config) {
		cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(layout))
		}
	}
}

func withServiceName(name string) LoggerOption {
	return func(cfg *zap.Config) {
		cfg.InitialFields = map[string]interface{}{
			"service": name,
		}
	}
}
