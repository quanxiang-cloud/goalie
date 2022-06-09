package config

import (
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"gopkg.in/yaml.v2"

	"io/ioutil"
)

// Conf 全局配置文件
var Conf *Config

// DefaultPath 默认配置路径
var DefaultPath = "./configs/config.yml"

// Config 配置文件
type Config struct {
	Port        string        `yaml:"port"`
	Model       string        `yaml:"model"`
	InternalNet client.Config `yaml:"internalNet"`
	Log         logger.Config `yaml:"log"`
	Mysql       mysql.Config  `yaml:"mysql"`
	Redis       redis.Config  `yaml:"redis"`
	Topic       string        `yaml:"topic"`
}

// NewConfig 获取配置配置
func NewConfig(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		return nil, err
	}

	return Conf, nil
}
