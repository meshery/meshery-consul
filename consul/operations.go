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
	"context"
	"fmt"
	"strings"

	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"

	"github.com/layer5io/meshery-adapter-library/adapter"
	opstatus "github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshery-consul/internal/config"
)

func (h *Consul) ApplyOperation(ctx context.Context, request adapter.OperationRequest) error {
	operations := make(adapter.Operations)
	err := h.Config.GetObject(adapter.OperationsKey, &operations)
	if err != nil {
		return err
	}

	status := opstatus.Deploying
	e := &adapter.Event{
		Operationid: request.OperationID,
		Summary:     "Deploying",
		Details:     "None",
	}

	if request.IsDeleteOperation {
		status = opstatus.Removing
		e.Summary = "Removing"
	}

	operation, ok := operations[request.OperationName]
	if !ok {
		e.Summary = "Error unknown operation name"
		err = adapter.ErrOpInvalid
		e.Details = err.Error()
		h.StreamErr(e, err)
		return err
	}

	opDesc := operation.Description
	kubeClient, err := mesherykube.New(h.KubeClient, h.RestConfig)
	if err != nil {
		e.Summary = fmt.Sprintf("Error while %s %s", status, opDesc)
		e.Details = err.Error()
		h.StreamErr(e, err)
		return err
	}

	switch request.OperationName {
	case config.Consul191DemoOperation: // Apply Helm chart operations
		if status, err = h.applyHelmChart(request, *operation, *kubeClient); err != nil {
			e.Summary = fmt.Sprintf("Error while %s %s", status, opDesc)
			e.Details = err.Error()
			h.StreamErr(e, err)
			return err
		}
	case config.CustomOperation, // Apply Kubernetes manifests operations
		config.Consul182DemoOperation,
		config.HTTPBinOperation,
		config.ImageHubOperation,
		config.BookInfoOperation:

		if status, err = h.applyManifests(request, *operation, *kubeClient); err != nil {
			e.Summary = fmt.Sprintf("Error while %s %s", status, opDesc)
			e.Details = err.Error()
			h.StreamErr(e, err)
			return err
		}

		e.Summary = fmt.Sprintf("%s %s successfully.", opDesc, status)
		e.Details = e.Summary

	default:
		h.StreamErr(e, adapter.ErrOpInvalid)
		return adapter.ErrOpInvalid
	}

	if !request.IsDeleteOperation && len(operation.Services) > 0 {
		for _, service := range operation.Services {
			svc := strings.TrimSpace(string(service))
			if len(svc) > 0 {
				h.Log.Info(fmt.Sprintf("Retreiving endpoint for service %s.", svc))

				endpoint, err1 := kubeClient.GetServiceEndpoint(ctx, svc, request.Namespace)
				if err1 != nil {
					h.StreamErr(&adapter.Event{
						Operationid: request.OperationID,
						Summary:     fmt.Sprintf("Unable to retrieve service endpoint for the service %s.", svc),
						Details:     err1.Error(),
					}, err1)
				} else {
					msg := fmt.Sprintf("%s Service endpoint %s at %s:%v", e.Summary, endpoint.Name, endpoint.Address, endpoint.Port)
					h.Log.Info(msg)
					e.Summary = msg
					e.Details = msg
				}
			}
		}
	}
	h.StreamInfo(e)

	return nil
}
