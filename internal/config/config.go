package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Bind   string `yaml:"bind"`
	Remote string `yaml:"remote"`
	Filter string `yaml:"filter-filename"`
	Log    string `yaml:"log-filename"`
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

func OpenLogFile(filename string) {
	if filename != "" {
		logfile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

		if err != nil {
			log.Fatalln("OpenLogfile: os.OpenFile: ", err)
		}

		log.SetOutput(logfile)
	}
}
