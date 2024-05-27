package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Conf struct {
	Url        string `yaml:"source_url"`
	Goroutines int    `yaml:"parallel"`
	Port       string `yaml:"port"`
	Dsn        string `yaml:"dsn"`
	CLimit     int    `yaml:"concurrency_limit"`
	RLimit     int64  `yaml:"rate_limit"`
	TokenTime  int    `yaml:"token_max_time"`
}

func (c *Conf) GetConf(path string) *Conf {

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
