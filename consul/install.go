package consul

import (
	"fmt"
	"path"

	"github.com/layer5io/gokit/errors"

	"github.com/layer5io/meshery-consul/internal/config"

	"github.com/layer5io/meshery-adapter-library/adapter"
)

type MeshInstance struct{}

func (m *MeshInstance) applyUsingManifest(request adapter.OperationRequest, operation *adapter.Operation, h *Handler) error {
	err := h.ApplyKubernetesManifest(request, *operation, map[string]string{"namespace": request.Namespace},
		path.Join("consul", "config_templates", operation.Properties[config.OperationTemplateNameKey]))
	return err
}

func (m *MeshInstance) getServicePorts(request adapter.OperationRequest, operation *adapter.Operation, h *Handler) ([]int64, error) {
	return h.GetServicePorts(operation.Properties[config.OperationServiceNameKey], request.Namespace)
}

func (h *Handler) applyUsingManifest(request adapter.OperationRequest, operation *adapter.Operation) (string, error) {
	status := "installing"

	if request.IsDeleteOperation {
		status = "removing"
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

	return "deployed", nil
}

func (h *Handler) getServicePorts(request adapter.OperationRequest, operation *adapter.Operation) (string, []int64, error) {
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
		error := errors.New("0000", fmt.Sprintf("Unable to retrieve information from the mesh: %s.", err.Error()))
		h.Log.Err(fmt.Sprintf("Retreiving port(s) for service %s failed.", svc), error.Error())
		return "", nil, error
	}

	var portMsg string
	if len(ports) == 1 {
		portMsg = fmt.Sprintf("The service %s is possibly available on port: %v.", svc, ports)
	} else if len(ports) > 1 {
		portMsg = fmt.Sprintf("The service %s is possibly available on one of the following ports: %v.", svc, ports)
	}

	return portMsg, ports, nil
}
