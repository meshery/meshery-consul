package consul

import (
	"github.com/layer5io/gokit/logger"
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/config"
)

type Handler struct {
	adapter.BaseHandler
}

func New(config config.Handler, log logger.Handler) adapter.Handler {
	return &Handler{
		adapter.BaseHandler{Config: config, Log: log},
	}
}
