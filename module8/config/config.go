// Package config
/*
 * @Author:       Kit-Hung
 * @Date:         2024/1/30 20:22
 * @Description： 总体的配置文件
 */
package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var Config = &GlobalConfig{
	Log: NewLogConfig(),
}

type GlobalConfig struct {
	Log *LogConfig `yaml:"log" json:"log"`
}

func InitGlobalConfig(filePath string) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(bytes, Config)
	if err != nil {
		log.Fatal(err)
	}
}
