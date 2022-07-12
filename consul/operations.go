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
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-consul/internal/config"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

func (h *Consul) ApplyOperation(ctx context.Context, request adapter.OperationRequest, hchan *chan interface{}) error {
	err := h.CreateKubeconfigs(request.K8sConfigs)
	if err != nil {
		return err
	}
	h.SetChannel(hchan)
	kubeconfigs := request.K8sConfigs
	operations := make(adapter.Operations)
	err = h.Config.GetObject(adapter.OperationsKey, &operations)
	if err != nil {
		return err
	}

	//status := opstatus.Deploying
	e := &adapter.Event{
		Operationid: request.OperationID,
		Summary:     "Deploying",
		Details:     "None",
	}

	if request.IsDeleteOperation {
		//status = opstatus.Removing
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

	switch request.OperationName {
	case config.ConsulOperation: // Apply Helm chart operations
		if status, err := h.applyHelmChart(request.IsDeleteOperation, operation.AdditionalProperties[config.HelmChartVersionKey], request.Namespace, kubeconfigs); err != nil {
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
		status, err := h.applyManifests(request, *operation, kubeconfigs)
		if err != nil {
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
	var errs []error
	var wg sync.WaitGroup
	for _, k8sconfig := range kubeconfigs {
		wg.Add(1)
		go func(k8sconfig string) {
			defer wg.Done()
			kClient, err := mesherykube.New([]byte(k8sconfig))
			if err != nil {
				errs = append(errs, err)
				return
			}
			if !request.IsDeleteOperation && len(operation.Services) > 0 {
				for _, service := range operation.Services {
					svc := strings.TrimSpace(string(service))
					if len(svc) > 0 {
						h.Log.Info(fmt.Sprintf("Retreiving endpoint for service %s.", svc))

						endpoint, err1 := mesherykube.GetServiceEndpoint(ctx, kClient.KubeClient, &mesherykube.ServiceOptions{
							Name:         svc,
							Namespace:    request.Namespace,
							APIServerURL: kClient.RestConfig.Host,
						})
						if err1 != nil {
							h.StreamErr(&adapter.Event{
								Operationid: request.OperationID,
								Summary:     fmt.Sprintf("Unable to retrieve service endpoint for the service %s.", svc),
								Details:     err1.Error(),
							}, err1)
						} else {
							external := "N/A"
							if endpoint.External != nil {
								external = fmt.Sprintf("%s:%v", endpoint.External.Address, endpoint.External.Port)
							}
							internal := "N/A"
							if endpoint.Internal != nil {
								internal = fmt.Sprintf("%s:%v", endpoint.Internal.Address, endpoint.Internal.Port)
							}
							msg := fmt.Sprintf("%s Service endpoints for service %s: internal=%s, external=%s", e.Summary, svc, internal, external)
							h.Log.Info(msg)
							e.Summary = msg
							e.Details = msg
						}
					}
				}
			}
		}(k8sconfig)
	}
	wg.Wait()

	h.StreamInfo(e)

	return nil
}
