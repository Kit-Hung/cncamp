// Package log
/*
 * @Author:       Kit-Hung
 * @Date:         2024/1/30 16:24
 * @Description： 日志相关
 */

package log

import (
	"github.com/Kit-Hung/cncamp/module10/config"
	"github.com/Kit-Hung/cncamp/module10/consts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var levelMap map[string]zapcore.Level

func init() {
	levelMap = make(map[string]zapcore.Level)
	levelMap[consts.LogLevelInfo] = zapcore.InfoLevel
	levelMap[consts.LogLevelDebug] = zapcore.DebugLevel

	productionConfig := getConfig()
	Logger = zap.Must(productionConfig.Build())
}

func getConfig() zap.Config {
	logConfig := config.Config.Log
	productionConfig := zap.NewProductionConfig()
	productionConfig.Level = zap.NewAtomicLevelAt(levelMap[logConfig.Level])
	productionConfig.OutputPaths = logConfig.OutputPaths
	productionConfig.EncoderConfig = getEncoderConfig()
	productionConfig.Encoding = "console"
	return productionConfig
}

func getEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	zapcore.NewConsoleEncoder(encoderConfig)
	return encoderConfig
}
