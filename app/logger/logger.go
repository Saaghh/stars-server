package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level string
}

func InitLogger(cfg Config) {
	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		panic(fmt.Errorf("zap.ParseAtomicLevel(%s): %w", cfg.Level, err))
	}

	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = zapcore.ISO8601TimeEncoder

	zap.ReplaceGlobals(zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(config), os.Stdout, level)))

	zap.L().Info("successful logger initialization")
}
