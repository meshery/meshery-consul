package config

import (
	"strings"

	smp "github.com/layer5io/service-mesh-performance/spec"
)

var (
	ConsulOperation = strings.ToLower(smp.ServiceMesh_CONSUL.Enum().String())
)
