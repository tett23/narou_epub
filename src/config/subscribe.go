package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	NCodes []string `yaml:"n_codes"`
}

func GetConfig() (*Config, error) {
	data, err := ioutil.ReadFile("./config/subscribe.yml")
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		panic(err)
	}

	return &config, err
}
