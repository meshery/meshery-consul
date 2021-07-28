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

const (
	CustomOperation         = "custom"
	Consul182DemoOperation  = "consul_182_demo"
	Consul191DemoOperation  = "consul_191_demo"
	BookInfoOperation       = "bookinfo"
	HTTPBinOperation        = "httpbin"
	ImageHubOperation       = "imagehub"
	SmiConformanceOperation = "smi_conformance"

	// Keys in AdditionalProperties of Operation
	HelmChartRepositoryKey = "helm_chart_repository"
	HelmChartChartKey      = "helm_chart_chart"
	HelmChartVersionKey    = "helm_chart_version"
	HelmChartValuesFileKey = "helm_chart_values_file"
)
