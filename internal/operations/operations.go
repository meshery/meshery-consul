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
				config.OperationDescriptionKey:  "Consul Connect: unsecured, 1 server, suitable for local exploration",
				config.OperationVersionKey:      "1.8.2",
				config.OperationTemplateNameKey: "consul.yaml",
			},
		},
		config.InstallBookInfoCommand: &adapter.Operation{
			Type: int32(meshes.OpCategory_SAMPLE_APPLICATION),
			Properties: map[string]string{
				config.OperationDescriptionKey:  "Istio Book Info Application",
				config.OperationVersionKey:      "",
				config.OperationTemplateNameKey: "bookinfo.yaml",
			},
		},
		config.InstallHTTPBinCommand: &adapter.Operation{
			Type: int32(meshes.OpCategory_SAMPLE_APPLICATION),
			Properties: map[string]string{
				config.OperationDescriptionKey:  "HTTPbin Application",
				config.OperationVersionKey:      "",
				config.OperationTemplateNameKey: "httpbin-consul.yaml",
			},
		},
		config.InstallImageHubCommand: &adapter.Operation{
			Type: int32(meshes.OpCategory_SAMPLE_APPLICATION),
			Properties: map[string]string{
				config.OperationDescriptionKey:  "Image Hub Application",
				config.OperationVersionKey:      "",
				config.OperationTemplateNameKey: "image-hub.yaml",
			},
		},
		config.CustomOpCommand: &adapter.Operation{
			Type: int32(meshes.OpCategory_CUSTOM),
			Properties: map[string]string{
				config.OperationDescriptionKey:  "Custom YAML",
				config.OperationVersionKey:      "",
				config.OperationTemplateNameKey: "image-hub.yaml",
			},
		},
	}
)
