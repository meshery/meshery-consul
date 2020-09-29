package config

import (
	"github.com/mgfeller/common-adapter-library/adapter"
	"github.com/mgfeller/common-adapter-library/config"
	"github.com/mgfeller/common-adapter-library/configprovider"
)

func New(name string, serverConfig map[string]string, meshSpec map[string]string, meshInstance map[string]string, providerConfig map[string]string, operations adapter.Operations) (config.Handler, error) {
	switch name {
	case "viper":
		return configprovider.NewViper(serverConfig, meshSpec, meshInstance, providerConfig, operations)
	}
	return nil, config.ErrEmptyConfig
}
