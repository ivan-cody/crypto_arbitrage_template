package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MempoolURL    string  `yaml:"mempool_url"`
	MinAmount     float64 `yaml:"min_amount"`
	Environment   string  `yaml:"environment"`
	SentryDSN     string  `yaml:"sentry_dsn"`
	EncryptionKey string  `yaml:"encryption_key"`
	DataDogConfig struct {
		APIKey string `yaml:"api_key"`
		AppKey string `yaml:"app_key"`
	} `yaml:"datadog"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
