package consul

import (
	"path"

	"github.com/layer5io/meshery-consul/internal/config"

	"github.com/mgfeller/common-adapter-library/adapter"
)

type MeshInstance struct{}

func (m *MeshInstance) applyConsulUsingManifest(request adapter.OperationRequest, operation *adapter.Operation, h *Handler) error {
	err := h.ApplyKubernetesManifest(request, *operation, map[string]string{},
		path.Join("consul", "config_templates", operation.Properties[config.OperationTemplateNameKey]))
	return err
}

func (h *Handler) applyConsulUsingManifest(request adapter.OperationRequest, operation *adapter.Operation) (string, error) {
	status := "installing" // TODO: should be type/enum defined in the common adapter library

	if request.IsDeleteOperation {
		status = "removing"
	}

	meshInstance := &MeshInstance{}

	err := h.Config.MeshInstance(meshInstance)
	if err != nil {
		return status, adapter.ErrMeshConfig(err)
	}

	h.Log.Info("Installing Consul")
	err = meshInstance.applyConsulUsingManifest(request, operation, h)
	if err != nil {
		h.Log.Err("Consul installation failed", adapter.ErrInstallMesh(err).Error())
		return status, adapter.ErrInstallMesh(err)
	}

	return "deployed", nil
}
