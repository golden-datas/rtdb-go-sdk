package config

import (
	"fmt"
	"os"

	"github.com/golden-datas/rtdb-go-sdk"
	"gopkg.in/yaml.v3"
)

// DBConfig 数据库配置
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Load 从YAML文件加载配置
func Load(path string) (*DBConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}

	var cfg DBConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file failed: %w", err)
	}

	// 验证必需字段
	if cfg.Host == "" {
		return nil, fmt.Errorf("config field 'host' is required")
	}
	if cfg.Port == 0 {
		return nil, fmt.Errorf("config field 'port' is required")
	}
	if cfg.Username == "" {
		return nil, fmt.Errorf("config field 'username' is required")
	}

	return &cfg, nil
}

// Connect 建立数据库连接
func (c *DBConfig) Connect() (*rtdb_api.RtdbConnect, error) {
	conn, err := rtdb_api.Login(c.Host, c.Port, c.Username, c.Password, rtdb_api.RtdbPrecisionSecond)
	if err != nil {
		return nil, fmt.Errorf("connect to database failed: %w", err)
	}
	return conn, nil
}

// GenerateTemplate 生成配置文件模板
func GenerateTemplate(path string) error {
	cfg := DBConfig{
		Host:     "127.0.0.1",
		Port:     6327,
		Username: "root",
		Password: "root",
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("marshal config failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config file failed: %w", err)
	}

	return nil
}
