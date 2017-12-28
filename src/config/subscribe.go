package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	NCodes            []string `yaml:"n_codes"`
	Email             string   `yaml:"email"`
	SendToKindleEmail string   `yaml:"send_to_kindle_email"`
	SMTPUserName      string   `yaml:"smtp_user_name"`
	SMTPPassword      string   `yaml:"smtp_password"`
	SMTPHost          string   `yaml:"smtp_host"`
	SMTPPort          int      `yaml:"smtp_port"`
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
