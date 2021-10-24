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

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"

	opstatus "github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

const (
	repo  = "https://helm.releases.hashicorp.com"
	chart = "consul"
)

func (h *Consul) installConsul(del bool, version, namespace string) (string, error) {
	h.Log.Debug(fmt.Sprintf("Requested install of version: %s", version))
	h.Log.Debug(fmt.Sprintf("Requested action is delete: %v", del))
	h.Log.Debug(fmt.Sprintf("Requested action is in namespace: %s", namespace))

	st := opstatus.Installing
	if del {
		st = opstatus.Removing
	}

	err := h.Config.GetObject(adapter.MeshSpecKey, h)
	if err != nil {
		return st, ErrMeshConfig(err)
	}

	_, err = h.applyHelmChart(del, version, namespace)
	if err != nil {
		return st, ErrApplyHelmChart(err)
	}

	st = opstatus.Installed
	if del {
		st = opstatus.Removed
	}

	return st, nil
}
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

// applyHelmChart installs or deletes an application packaged as Helm chart.
//
// Information about the Helm chart is specified in operation additional properties (keys defined in keys.go):
// the repository, the chart, the version, and the name of a values-file.
// The chart is the only required value, defaults are handled by ApplyHelmChart from meshkit.
// The values-file is expected in the consul/config_templates folder. It corresponds to a file specified
// by the -f (--values) flag of the Helm CLI.
func (h *Consul) applyHelmChart(del bool, version string, ns string) (string, error) {
	status := opstatus.Installing
	var act mesherykube.HelmChartAction
	if del {
		status = opstatus.Removing
		act = mesherykube.UNINSTALL
	} else {
		act = mesherykube.INSTALL
	}

	err := h.MesheryKubeclient.ApplyHelmChart(mesherykube.ApplyHelmChartConfig{
		Namespace:       ns,
		CreateNamespace: true,
		Action:          act,
		ChartLocation: mesherykube.HelmChartLocation{
			Repository: repo,
			Chart:      chart,
			Version:    version,
		},
	})
	if err != nil {
		return status, ErrApplyOperation(err)
	}

	return opstatus.Deployed, nil
}
