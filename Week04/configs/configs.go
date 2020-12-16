package configs

import (
	"io/ioutil"

	"github.com/ghodss/yaml"

	xerrors "github.com/pkg/errors"
)

type Config struct {
	DataBaseAddress string `yaml:"data_base_address"`
	TableName       string `yaml:"table_name"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Init() error {
	file, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return xerrors.Wrapf(err, "failed reading config file")
	}

	return xerrors.Wrapf(yaml.Unmarshal(file, c), "failed unmarshaling config file")
}
