package app

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host string   `yaml:"host" json:"host" env:"SERVER_HOST"`
	Port string   `yaml:"port"`
	DB   Database `yaml:"database"` // в структуре yaml файла есть поле database
}

type Database struct {
	DSN string `yaml:"dsn"` // "postgres://user:password@localhost:5432/dbname"
}

func NewConfig(configPath string) (*Config, error) {
	var config = new(Config)
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
