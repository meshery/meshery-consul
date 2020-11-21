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
	opstatus "github.com/layer5io/meshery-adapter-library/status"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

func (h *Consul) applyManifests(request adapter.OperationRequest, operation adapter.Operation, kubeClient mesherykube.Client) (string, error) {
	status := opstatus.Installing

	if request.IsDeleteOperation {
		status = opstatus.Removing
	}

	h.Log.Info(fmt.Sprintf("%s %s", status, operation.Description))

	for _, template := range operation.Templates {
		p := path.Join("consul", "config_templates", string(template))
		m, err := ioutil.ReadFile(p)
		if err != nil {
			return status, ErrApplyOperation(err)
		}
		err = kubeClient.ApplyManifest(m, mesherykube.ApplyOptions{
			Namespace: request.Namespace,
			Update:    true,
			Delete:    request.IsDeleteOperation,
		})
		if err != nil {
			return status, ErrApplyOperation(err)
		}
	}

	return opstatus.Deployed, nil
}
