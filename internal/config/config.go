package config

import (
	"github.com/mgfeller/common-adapter-library/config"
)

func New(name string, serverConfig map[string]string, meshConfig map[string]string, providerConfig map[string]string) (config.Handler, error) {
	switch name {
	case "viper":
		return NewViper(serverConfig, meshConfig, providerConfig)
	}
	return nil, config.ErrEmptyConfig
}
