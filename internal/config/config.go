package config

import (
	"github.com/layer5io/meshery-adapter-library/config"
	"github.com/layer5io/meshery-adapter-library/config/provider"
)

func New(options provider.Options) (config.Handler, error) {
	return provider.NewViper(options)
}
