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
		go func(hh *ConsulAdapter, ee *adapter.Event) {
			if status, err := hh.installConsul(doDelete); err != nil {
				ee.Summary = fmt.Sprintf("Error while %s Consul service mesh", status)
				ee.Details = err.Error()
				hh.StreamErr(ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("Consul service mesh %s successfully", status)
			ee.Details = fmt.Sprintf("The Consul service mesh is now %s.", status)
			hh.StreamInfo(ee)
		}(h, e)
		h.StreamErr(e, adapter.ErrOpInvalid)
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
