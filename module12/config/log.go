// Package config
/*
 * @Author:       Kit-Hung
 * @Date:         2024/1/30 17:04
 * @Description： 日志相关配置
 */
package config

import (
	"github.com/Kit-Hung/cncamp/module12/consts"
)

type LogConfig struct {
	Level       string   `json:"level" yaml:"level"`
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`
}

func NewLogConfig() *LogConfig {
	return &LogConfig{
		Level:       consts.LogLevelInfo,
		OutputPaths: []string{"stdout"},
	}
}
