package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AppPort      string `yaml:"port"`
	DbConnString string `yaml:"db_conn_string"`
}

func Parse(confPath string) (*Config, error) {
	filename, err := filepath.Abs(confPath)
	if err != nil {
		return nil, fmt.Errorf("can't get config path: %s", err.Error())
	}
	yamlConf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can't read conf: %s", err.Error())
	}

	var config Config
	err = yaml.Unmarshal(yamlConf, &config)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshall conf: %s", err.Error())
	}

	return &config, nil
}
