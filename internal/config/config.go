package config

import (
	"github.com/mgfeller/common-adapter-library/adapter"
	"github.com/mgfeller/common-adapter-library/config"
	"github.com/mgfeller/common-adapter-library/config_provider"
)

func New(name string, serverConfig map[string]string, meshConfig map[string]string, providerConfig map[string]string, operations adapter.Operations) (config.Handler, error) {
	switch name {
	case "viper":
		return config_provider.NewViper(serverConfig, meshConfig, providerConfig, operations)
	}
	return nil, config.ErrEmptyConfig
}
