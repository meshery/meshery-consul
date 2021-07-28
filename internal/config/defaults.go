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

package config

import (
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	"github.com/layer5io/meshery-adapter-library/status"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var (
	ServerDefaults = map[string]string{
		"name":     smp.ServiceMesh_CONSUL.Enum().String(),
		"type":     "adapter",
		"port":     "10002",
		"traceurl": "none",
	}

	MeshSpecDefaults = map[string]string{
		"name":    smp.ServiceMesh_CONSUL.Enum().String(),
		"status":  status.NotInstalled,
		"version": "1.8.2",
		"type":    smp.ServiceMesh_CONSUL.Enum().String(),
	}

	ViperDefaults = map[string]string{
		"filepath": ConfigRootPath,
		"filename": "consul",
		"filetype": "yaml",
	}

	KubeConfigDefaults = map[string]string{
		configprovider.FilePath: ConfigRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "kubeconfig",
	}
)
