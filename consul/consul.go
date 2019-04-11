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

import (
	"context"
	"time"

	"github.com/layer5io/meshery-consul/meshes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (iClient *ConsulClient) CreateMeshInstance(_ context.Context, k8sReq *meshes.CreateMeshInstanceRequest) (*meshes.CreateMeshInstanceResponse, error) {
	iClient.eventChan = make(chan *meshes.EventsResponse, 100)
	return &meshes.CreateMeshInstanceResponse{}, nil
}

// MeshName just returns the name of the mesh the client is representing
func (iClient *ConsulClient) MeshName(context.Context, *meshes.MeshNameRequest) (*meshes.MeshNameResponse, error) {
	return &meshes.MeshNameResponse{Name: "Consul"}, nil
}

// ApplyRule is a method invoked to apply a particular operation on the mesh in a namespace
func (iClient *ConsulClient) ApplyOperation(ctx context.Context, arReq *meshes.ApplyRuleRequest) (*meshes.ApplyRuleResponse, error) {
	return &meshes.ApplyRuleResponse{}, nil
}

// SupportedOperations - returns a list of supported operations on the mesh
func (iClient *ConsulClient) SupportedOperations(context.Context, *meshes.SupportedOperationsRequest) (*meshes.SupportedOperationsResponse, error) {
	result := map[string]string{}
	for key, op := range supportedOps {
		result[key] = op.name
	}
	return &meshes.SupportedOperationsResponse{
		Ops: result,
	}, nil
}

// StreamEvents - streams generated/collected events to the client
func (iClient *ConsulClient) StreamEvents(in *meshes.EventsRequest, stream meshes.MeshService_StreamEventsServer) error {
	for {
		select {
		case event := <-iClient.eventChan:
			logrus.Debugf("sending event: %+#v", event)
			if err := stream.Send(event); err != nil {
				err = errors.Wrapf(err, "unable to send event")

				// to prevent loosing the event, will re-add to the channel
				go func() {
					iClient.eventChan <- event
				}()
				logrus.Error(err)
				return err
			}
		default:
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}
