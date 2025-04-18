package config

type Config struct {
	Chat   Chat   `yaml:"chat"`
	Logger Logger `yaml:"logger"`
}

type Chat struct {
	Http    Http   `yaml:"http"`
	LogFile string `yaml:"log-file"`
}

type Http struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Cors Cors   `yaml:"cors"`
}

type Cors struct {
	Origins          []string `yaml:"origins"`
	Headers          []string `yaml:"headers"`
	AllowCredentials bool     `yaml:"allow-credentials"`
}

type Logger struct {
	MaxSizeMB  int `yaml:"max-size-mb"`
	MaxBackups int `yaml:"max-backups"`
	MaxAgeDays int `yaml:"max-age-days"`
}
