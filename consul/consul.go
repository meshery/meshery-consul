package consul

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/config"
	"github.com/layer5io/meshkit/logger"
)

type Handler struct {
	adapter.Adapter
}

func New(config config.Handler, log logger.Handler) adapter.Handler {
	return &Handler{
		adapter.Adapter{Config: config, Log: log},
	}
}
