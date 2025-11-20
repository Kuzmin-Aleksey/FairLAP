package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Debug            bool             `json:"debug" yaml:"debug" env:"DEBUG" envDefault:"false"`
	Http             *HttpConfig      `json:"http" yaml:"http"`
	MySQL            *MySQLConfig     `json:"mysql" yaml:"mysql"`
	YoloModel        *YoloModelConfig `json:"yolo_model" yaml:"yolo_model"`
	ImagesPath       string           `json:"images_path" yaml:"images_path"`
	DefaultLapConfig map[string]int   `json:"default_lap_config" yaml:"default_lap_config"`
}

type HttpConfig struct {
	Host             string `json:"host" yaml:"host"`
	ReadTimeoutSec   int    `json:"read_timeout_sec" yaml:"read_timeout_sec"`
	HandleTimeoutSec int    `json:"handle_timeout_sec" yaml:"handle_timeout_sec"`
	WriteTimeoutSec  int    `json:"write_timeout_sec" yaml:"write_timeout_sec"`
	SSLKeyPath       string `json:"ssl_key_path" yaml:"ssl_key_path"`
	SSLCertPath      string `json:"ssl_cert_path" yaml:"ssl_cert_path"`
}

type MySQLConfig struct {
	Host              string `json:"host" yaml:"host"`
	Port              int    `json:"port" yaml:"port"`
	User              string `json:"user" yaml:"user"`
	Password          string `json:"password" yaml:"password"`
	Schema            string `json:"schema" yaml:"schema"`
	ConnectTimeoutSec int    `json:"connect_timeout_sec" yaml:"connect_timeout_sec"`
}

type YoloModelConfig struct {
	Model       string `json:"model" yaml:"model"`
	ModelConfig string `json:"model_config" yaml:"model_config"`
}

func ReadConfig(path string, dotenv ...string) (*Config, error) {
	if err := godotenv.Load(dotenv...); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	cfg := new(Config)
	cfg.DefaultLapConfig = make(map[string]int)

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
