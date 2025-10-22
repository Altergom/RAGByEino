package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type Config struct {
	TargetURL         string        `yaml:"TargetURL"`
	CheckInterval     time.Duration `yaml:"CheckInterval"`
	Timeout           time.Duration `yaml:"Timeout"`
	Headless          bool          `yaml:"Headless"`
	BrowserType       string        `yaml:"BrowserType"`
	ChromeDriverPath  string        `yaml:"ChromeDriverPath"`
	ChromeDriverPort  int           `yaml:"ChromeDriverPort"`
	FirefoxDriverPath string        `yaml:"FirefoxDriverPath"`
}

// defaultConfig 默认配置
func defaultConfig() Config {
	return Config{
		TargetURL:         "https://baidu.com",
		CheckInterval:     5 * time.Second,
		Timeout:           60 * time.Second,
		Headless:          false,
		BrowserType:       "chrome", // 默认使用Chrome
		ChromeDriverPath:  "./chrome/chromedriver-win64/chromedriver",
		ChromeDriverPort:  4444,
		FirefoxDriverPath: "./firefox/geckodriver-v0.36.0-win64/geckodriver",
	}
}

func InitConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件出错: %w", err)
	}

	var cfg *Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件出错: %w", err)
	}

	defaultCfg := defaultConfig()

	if cfg.TargetURL == "" {
		cfg.TargetURL = defaultCfg.TargetURL
	}
	if cfg.CheckInterval == 0 {
		cfg.CheckInterval = defaultCfg.CheckInterval
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultCfg.Timeout
	}
	if cfg.Headless {
		cfg.Headless = defaultCfg.Headless
	}
	if cfg.BrowserType == "" {
		cfg.BrowserType = defaultCfg.BrowserType
	}
	if cfg.ChromeDriverPath == "" {
		cfg.ChromeDriverPath = defaultCfg.ChromeDriverPath
	}
	if cfg.ChromeDriverPort == 0 {
		cfg.ChromeDriverPort = defaultCfg.ChromeDriverPort
	}
	if cfg.FirefoxDriverPath == "" {
		cfg.FirefoxDriverPath = defaultCfg.FirefoxDriverPath
	}
	return cfg, nil
}
