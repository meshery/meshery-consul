package config

import (
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var (
	ConsulOperation = strings.ToLower(smp.ServiceMesh_CONSUL.Enum().String())
)

func getOperations(dev adapter.Operations) adapter.Operations {
	versions, _ := GetLatestReleases(3)
	var versionNames []adapter.Version
	for _, v := range versions {
		versionNames = append(versionNames, v.Name)
	}
	dev[ConsulOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_INSTALL),
		Description: "Consul",
		Versions:    versionNames,
		Templates: []adapter.Template{
			"templates/consul.yaml",
		},
		AdditionalProperties: map[string]string{},
	}

	return dev
}
