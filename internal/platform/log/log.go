package log

import (
	"strings"

	"biz/internal/platform/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.LogConfig, profile string) (*zap.Logger, error) {
	zcfg := zap.NewProductionConfig()
	if strings.EqualFold(profile, "dev") {
		zcfg = zap.NewDevelopmentConfig()
	}
	if strings.EqualFold(cfg.Format, "console") {
		zcfg.Encoding = "console"
	} else {
		zcfg.Encoding = "json"
	}
	lvl := zapcore.InfoLevel
	_ = lvl.Set(strings.ToLower(cfg.Level))
	zcfg.Level = zap.NewAtomicLevelAt(lvl)
	return zcfg.Build()
}
