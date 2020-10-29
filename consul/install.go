package consul

import (
	"fmt"
	"path"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-consul/internal/config"
	opstatus "github.com/layer5io/meshery-consul/internal/status"
)

type MeshInstance struct{}

func (m *MeshInstance) applyUsingManifest(request adapter.OperationRequest, operation adapter.Operation, h *Handler) error {
	err := h.ApplyKubernetesManifest(request, operation, map[string]string{"namespace": request.Namespace},
		path.Join("consul", "config_templates", operation.Properties[config.OperationTemplateNameKey]))
	return err
}

func (m *MeshInstance) getServicePorts(request adapter.OperationRequest, operation adapter.Operation, h *Handler) ([]int64, error) {
	return h.GetServicePorts(operation.Properties[config.OperationServiceNameKey], request.Namespace)
}

func (h *Handler) applyUsingManifest(request adapter.OperationRequest, operation adapter.Operation) (string, error) {
	status := opstatus.Installing

	if request.IsDeleteOperation {
		status = opstatus.Removing
	}

	meshInstance := &MeshInstance{}

	err := h.Config.MeshInstance(meshInstance)
	if err != nil {
		return status, adapter.ErrMeshConfig(err)
	}

	h.Log.Info(fmt.Sprintf("%s %s", status, operation.Properties[config.OperationDescriptionKey]))
	err = meshInstance.applyUsingManifest(request, operation, h)
	if err != nil {
		h.Log.Err(fmt.Sprintf("%s %s failed", status, operation.Properties[config.OperationDescriptionKey]), adapter.ErrInstallMesh(err).Error())
		return status, adapter.ErrInstallMesh(err)
	}

	return opstatus.Deployed, nil
}

func (h *Handler) getServicePorts(request adapter.OperationRequest, operation adapter.Operation) (string, []int64, error) {
	meshInstance := &MeshInstance{}

	err := h.Config.MeshInstance(meshInstance)
	if err != nil {
		return "", nil, adapter.ErrMeshConfig(err)
	}

	svc := operation.Properties[config.OperationServiceNameKey]

	h.Log.Info(fmt.Sprintf("Retreiving service ports for service %s.", svc))

	var ports []int64
	ports, err = meshInstance.getServicePorts(request, operation, h)
	if err != nil {
		err2 := ErrGetInfo(err)
		h.Log.Err(fmt.Sprintf("Retreiving port(s) for service %s failed.", svc), err2.Error())
		return "", nil, err2
	}

	var portMsg string
	if len(ports) == 1 {
		portMsg = fmt.Sprintf("The service %s is possibly available on port: %v.", svc, ports)
	} else if len(ports) > 1 {
		portMsg = fmt.Sprintf("The service %s is possibly available on one of the following ports: %v.", svc, ports)
	}

	return portMsg, ports, nil
}
