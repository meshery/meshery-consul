package consul

import (
	"context"
	"errors"
	"fmt"

	"github.com/layer5io/meshery-consul/internal/config"
	"github.com/mgfeller/common-adapter-library/adapter"
)

func (h *Handler) ApplyOperation(ctx context.Context, request adapter.OperationRequest) error {
	operations := make(adapter.Operations)
	err := h.Config.Operations(&operations)
	if err != nil {
		return err
	}

	status := "deploying"
	e := &adapter.Event{
		Operationid: request.OperationID,
		Summary:     "Deploying",
		Details:     "None",
	}

	if request.IsDeleteOperation {
		status = "removing"
		e.Summary = "Removing"
	}

	if err := h.CreateNamespace(request.IsDeleteOperation, request.Namespace); err != nil {
		e.Summary = "Error while creating namespace"
		e.Details = err.Error()
		h.StreamErr(e, err)
		return err
	}

	operation, ok := operations[request.OperationName]
	if !ok {
		e.Summary = "Error unknown operation name"
		err = errors.New(e.Summary)
		e.Details = err.Error()
		h.StreamErr(e, err)
		return err
	}

	switch request.OperationName {
	case config.CustomOpCommand:
		h.StreamErr(e, adapter.ErrOpInvalid)
	case config.InstallConsulCommand:
		if status, err := h.applyConsulUsingManifest(request, operation); err != nil {
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
