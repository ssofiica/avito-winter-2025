package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Host              string        `yaml:"host"`
	Port              string        `yaml:"port"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"`
}

func Load() *Config {
	config := &Config{}
	buf, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatal("Error reading application configuration: ", err.Error())
	}

	err = yaml.Unmarshal(buf, config)
	if err != nil {
		log.Fatal("Error creating configuration object: ", err.Error())
	}

	fmt.Println("Reading configuration successful")
	return config
}
