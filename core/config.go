package core

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Hoods []*Hood `yaml:"hood"`
}

func Load(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(data, c)
	return c, err
}