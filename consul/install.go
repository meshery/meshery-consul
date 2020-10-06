package consul

import (
	"fmt"
	"path"

	"github.com/layer5io/meshery-consul/internal/config"

	"github.com/layer5io/meshery-adapter-library/adapter"
)

type MeshInstance struct{}

func (m *MeshInstance) applyUsingManifest(request adapter.OperationRequest, operation *adapter.Operation, h *Handler) error {
	err := h.ApplyKubernetesManifest(request, *operation, map[string]string{"namespace": request.Namespace},
		path.Join("consul", "config_templates", operation.Properties[config.OperationTemplateNameKey]))
	return err
}

func (h *Handler) applyUsingManifest(request adapter.OperationRequest, operation *adapter.Operation) (string, error) {
	status := "installing" // TODO: should be enum defined in the common adapter library

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
		h.Log.Err(fmt.Sprintf("Error: %s %s failed", status, operation.Properties[config.OperationDescriptionKey]), adapter.ErrInstallMesh(err).Error())
		return status, adapter.ErrInstallMesh(err)
	}

	return "deployed", nil
}
