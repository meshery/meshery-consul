package consul

import "github.com/mgfeller/common-adapter-library/adapter"

type MeshInstance struct{}

func (m *MeshInstance) applyConsulUsingManifest(doDelete bool) error {

	return nil
}

func (h *ConsulAdapter) installConsul(doDelete bool) (string, error) {
	status := "installing"

	if doDelete {
		status = "removing"
	}

	meshInstance := &MeshInstance{}

	err := h.Config.MeshInstance(meshInstance)
	if err != nil {
		return status, adapter.ErrMeshConfig(err)
	}

	h.Log.Info("Installing Consul")
	err = meshInstance.applyConsulUsingManifest(doDelete)
	if err != nil {
		h.Log.Err("Consul installation failed", adapter.ErrInstallMesh(err).Error())
		return status, adapter.ErrInstallMesh(err)
	}

	return "deployed", nil
}
