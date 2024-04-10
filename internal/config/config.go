package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Conf struct {
	Url string `yaml:"source_url"`
	Bd  string `yaml:"db_file"`
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
