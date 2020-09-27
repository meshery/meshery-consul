package consul

import (
	"context"
	"fmt"
	"github.com/layer5io/meshery-consul/internal/config"
	"github.com/mgfeller/common-adapter-library/adapter"
)

func (h *ConsulAdapter) ApplyOperation(ctx context.Context, op string, id string, doDelete bool) error {

	operations := make(adapter.Operations, 0)
	err := h.Config.Operations(&operations)
	if err != nil {
		return err
	}

	status := "deploying"
	e := &adapter.Event{
		Operationid: id,
		Summary:     "Deploying",
		Details:     "None",
	}

	if doDelete {
		status = "removing"
		e.Summary = "Removing"
	}

	switch op {
	case config.CustomOpCommand:
		h.StreamErr(e, adapter.ErrOpInvalid)
	case config.InstallConsulCommand:
		if status, err := h.installConsul(doDelete); err != nil {
			e.Summary = fmt.Sprintf("Error while %s Consul service mesh", status)
			e.Details = err.Error()
			h.StreamErr(e, err)
			return err
		}
		e.Summary = fmt.Sprintf("Consul service mesh %s successfully", status)
		e.Details = fmt.Sprintf("The Consul service mesh is now %s.", status)
		h.StreamInfo(e)
	case config.InstallHTTPBinCommand:
		h.StreamErr(e, adapter.ErrOpInvalid)
	case config.InstallImageHubCommand:
		h.StreamErr(e, adapter.ErrOpInvalid)
	case config.InstallBookInfoCommand:
		h.StreamErr(e, adapter.ErrOpInvalid)
	default:
		h.StreamErr(e, adapter.ErrOpInvalid)
	}
	return nil
}
