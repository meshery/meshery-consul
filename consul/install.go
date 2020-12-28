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

package consul

import (
	"fmt"
	"io/ioutil"
	"path"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
	opstatus "github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshery-consul/internal/config"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

func (h *Consul) applyManifests(request adapter.OperationRequest, operation adapter.Operation, kubeClient mesherykube.Client) (string, error) {
	status := opstatus.Installing

	if request.IsDeleteOperation {
		status = opstatus.Removing
	}

	h.Log.Info(fmt.Sprintf("%s %s", status, operation.Description))

	if operation.Type == int32(meshes.OpCategory_CUSTOM) {
		err := kubeClient.ApplyManifest([]byte(request.CustomBody), mesherykube.ApplyOptions{
			Namespace: request.Namespace,
			Update:    true,
			Delete:    request.IsDeleteOperation,
		})
		if err != nil {
			return status, ErrApplyOperation(err)
		}
	} else {
		for _, template := range operation.Templates {
			p := path.Join("consul", "config_templates", string(template))
			tpl, err := ioutil.ReadFile(p)
			if err != nil {
				return status, ErrApplyOperation(err)
			}
			merged, err := utils.MergeToTemplate(tpl, map[string]string{"namespace": request.Namespace})
			if err != nil {
				return status, ErrApplyOperation(err)
			}
			err = kubeClient.ApplyManifest(merged, mesherykube.ApplyOptions{
				Namespace: request.Namespace,
				Update:    true,
				Delete:    request.IsDeleteOperation,
			})
			if err != nil {
				return status, ErrApplyOperation(err)
			}
		}
	}
	return opstatus.Deployed, nil
}

func (h *Consul) applyHelmChart(request adapter.OperationRequest, operation adapter.Operation, kubeClient mesherykube.Client) (string, error) {
	status := opstatus.Installing

	if request.IsDeleteOperation {
		status = opstatus.Removing
	}

	h.Log.Info(fmt.Sprintf("%s %s", status, operation.Description))

	p := getter.All(cli.New())
	valueOpts := &values.Options{}
	if valuesFile, ok := operation.AdditionalProperties[config.HelmChartValuesFileKey]; ok {
		valueOpts.ValueFiles = []string{path.Join("consul", "config_templates", valuesFile)}
	}
	values, err := valueOpts.MergeValues(p)
	if err != nil {
		return status, ErrApplyOperation(err)
	}

	err = kubeClient.ApplyHelmChart(mesherykube.ApplyHelmChartConfig{
		Namespace:       request.Namespace,
		CreateNamespace: true,
		Delete:          request.IsDeleteOperation,
		ChartLocation: mesherykube.HelmChartLocation{
			Repository: operation.AdditionalProperties[config.HelmChartRepositoryKey],
			Chart:      operation.AdditionalProperties[config.HelmChartChartKey],
			Version:    operation.AdditionalProperties[config.HelmChartVersionKey],
		},
		OverrideValues: values,
	})
	if err != nil {
		return status, ErrApplyOperation(err)
	}

	return opstatus.Deployed, nil
}
