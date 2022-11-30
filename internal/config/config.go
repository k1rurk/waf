package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	Bind     string `yaml:"bind"`
	Remote   string `yaml:"remote"`
	Filename string `yaml:"filter-filename"`
}

func ReadConfigFile() *Config {
	config := new(Config)

	yamlFile, err := ioutil.ReadFile("configs/config.yaml")
	if err != nil {
		log.Printf("Reading file error %v\n", err)
	}
	err = yaml.Unmarshal(yamlFile, config)

	if err != nil {
		log.Fatalf("Unmarshal: %v\n", err)
	}

	return config
}
