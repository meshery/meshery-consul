// Copyright 2020 Layer5, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package operations

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
	"github.com/layer5io/meshery-consul/internal/config"
)

var (
	Operations = adapter.Operations{
		config.Consul182DemoOperation: &adapter.Operation{
			Type:        int32(meshes.OpCategory_INSTALL),
			Description: "Consul Connect v1.8.2: unsecured, 1 server, suitable for local exploration",
			Versions:    []adapter.Version{"1.8.2"},
			Templates:   []adapter.Template{"consul-1.8.2-demo.yaml"},
			Services:    []adapter.Service{"consul-consul-ui"},
		},
		config.Consul191DemoOperation: &adapter.Operation{
			Type:        int32(meshes.OpCategory_INSTALL),
			Description: "Consul Connect v1.9.1: unsecured, 1 server, suitable for local exploration",
			Versions:    []adapter.Version{"1.9.1"},
			Templates:   []adapter.Template{},
			Services:    []adapter.Service{"consul-consul-ui"},
			AdditionalProperties: map[string]string{
				config.HelmChartRepositoryKey: "https://helm.releases.hashicorp.com",
				config.HelmChartChartKey:      "consul",
				config.HelmChartVersionKey:    "0.28.0",
				config.HelmChartValuesFileKey: "consul-values-1.9.1-demo.yaml",
			},
		},
		config.BookInfoOperation: &adapter.Operation{
			Type:        int32(meshes.OpCategory_SAMPLE_APPLICATION),
			Description: "Istio Book Info Application",
			Versions:    []adapter.Version{},
			Templates:   []adapter.Template{"bookinfo.yaml"},
			Services:    []adapter.Service{"productpage"},
		},
		config.HTTPBinOperation: &adapter.Operation{
			Type:        int32(meshes.OpCategory_SAMPLE_APPLICATION),
			Description: "HTTPbin Application",
			Versions:    []adapter.Version{},
			Templates:   []adapter.Template{"httpbin-consul.yaml"},
			Services:    []adapter.Service{"httpbin"},
		},
		config.ImageHubOperation: &adapter.Operation{
			Type:        int32(meshes.OpCategory_SAMPLE_APPLICATION),
			Description: "Image Hub Application",
			Versions:    []adapter.Version{},
			Templates:   []adapter.Template{"image-hub.yaml"},
			Services:    []adapter.Service{"ingess"},
		},
		config.CustomOperation: &adapter.Operation{
			Type:        int32(meshes.OpCategory_CUSTOM),
			Description: "Custom YAML",
			Versions:    []adapter.Version{},
			Templates:   []adapter.Template{},
		},
		config.SmiConformanceOperation: &adapter.Operation{
			Type:        int32(meshes.OpCategory_VALIDATE),
			Description: "SMI Conformance",
		},
	}
)
