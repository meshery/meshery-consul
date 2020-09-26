package consul

import (
	"github.com/layer5io/gokit/logger"
	"github.com/mgfeller/common-adapter-library/adapter"
	"github.com/mgfeller/common-adapter-library/config"
)

type ConsulAdapter struct {
	adapter.BaseAdapter
}

func New(config config.Handler, log logger.Handler) adapter.Handler {
	return &ConsulAdapter{
		adapter.BaseAdapter{Config: config, Log: log},
	}
}
