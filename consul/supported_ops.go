// Copyright 2019 Layer5.io
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

package consul

import "github.com/layer5io/meshery-consul/meshes"

type supportedOperation struct {
	// a friendly name
	name string
	// the template file name
	templateName string

	opType meshes.OpCategory
}

const (
	customOpCommand        = "custom"
	installConsulCommand   = "consul_install"
	installBookInfoCommand = "install_book_info"
	installHTTPBinCommand  = "install_http_bin"
)

var supportedOps = map[string]supportedOperation{
	customOpCommand: {
		name:   "Custom YAML",
		opType: meshes.OpCategory_CUSTOM,
	},
	installConsulCommand: {
		name:         "Latest version of Consul with sidecar injector",
		templateName: "consul.yaml",
		opType:       meshes.OpCategory_INSTALL,
	},
	installBookInfoCommand: {
		name:         "Istio Book Info Application",
		templateName: "bookinfo.yaml",
		opType:       meshes.OpCategory_SAMPLE_APPLICATION,
	},
	installHTTPBinCommand: {
		name:         "HTTPbin Application",
		templateName: "httpbin-consul.yaml",
		opType:       meshes.OpCategory_SAMPLE_APPLICATION,
	},
}
