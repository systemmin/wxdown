// config/config.go

package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Port        string `yaml:"port"`
	Path        string `yaml:"path"`
	Browser     bool   `yaml:"browser"`
	Wkhtmltopdf struct {
		Enable bool   `yaml:"enable"`
		Path   string `yaml:"path"`
	} `yaml:"wkhtmltopdf"`
	Thread struct {
		Html  int `yaml:"html"`
		Image int `yaml:"image"`
	} `yaml:"thread"`
	Base64 bool `yaml:"base64"`
	Https  bool `yaml:"https"`
	Auth   struct {
		Enable bool     `yaml:"enable"`
		Users  []string `yaml:"users"`
	} `yaml:"auth"`
}

// LoadConfig 方法用于读取和解析 YAML 配置文件
func LoadConfig(base string) *Config {
	// 读取 YAML 配置文件
	file, err := os.ReadFile(filepath.Join(base, "config", "config.yaml"))
	if err != nil {
		file, err = os.ReadFile(filepath.Join(base, "config.yaml"))
		if err != nil {
			log.Fatalf("读取配置文件 config.yml 失败: %s", err)
		}
	}
	// 解析 YAML 数据
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Failed to parse YAML data: %s", err)
	}
	return &config
}
