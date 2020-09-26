package operations

import (
	"github.com/layer5io/meshery-consul/internal/config"
	"github.com/mgfeller/common-adapter-library/adapter"
	"github.com/mgfeller/common-adapter-library/meshes"
)

var (
	Operations = adapter.Operations{
		config.InstallConsulCommand: &adapter.Operation{
			Type: int32(meshes.OpCategory_INSTALL),
			Properties: map[string]string{
				"description":  "Consul Connect: unsecured, 1 server, suitable for local exploration",
				"version":      "1.8.2",
				"templateName": "consul.yaml",
			},
		},
		config.InstallBookInfoCommand: &adapter.Operation{
			Type: int32(meshes.OpCategory_SAMPLE_APPLICATION),
			Properties: map[string]string{
				"description":  "Istio Book Info Application",
				"version":      "",
				"templateName": "bookinfo.yaml",
			},
		},
		config.InstallHTTPBinCommand: &adapter.Operation{
			Type: int32(meshes.OpCategory_SAMPLE_APPLICATION),
			Properties: map[string]string{
				"description":  "HTTPbin Application",
				"version":      "",
				"templateName": "httpbin-consul.yaml",
			},
		},
		config.InstallImageHubCommand: &adapter.Operation{
			Type: int32(meshes.OpCategory_SAMPLE_APPLICATION),
			Properties: map[string]string{
				"description":  "Image Hub Application",
				"version":      "",
				"templateName": "image-hub.yaml",
			},
		},
		config.CustomOpCommand: &adapter.Operation{
			Type: int32(meshes.OpCategory_CUSTOM),
			Properties: map[string]string{
				"description":  "Custom YAML",
				"version":      "",
				"templateName": "image-hub.yaml",
			},
		},
	}
)
