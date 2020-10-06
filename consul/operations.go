package consul

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-consul/internal/config"
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
		err = adapter.ErrOpInvalid
		e.Details = err.Error()
		h.StreamErr(e, err)
		return err
	}

	switch request.OperationName {
	case config.CustomOpCommand, // TODO: implement custom op and test for it
		config.InstallConsulCommand,
		config.InstallHTTPBinCommand,
		config.InstallImageHubCommand,
		config.InstallBookInfoCommand:
		if status, err := h.applyUsingManifest(request, operation); err != nil {
			e.Summary = fmt.Sprintf("Error while %s %s", status, operation.Properties[config.OperationDescriptionKey])
			e.Details = err.Error()
			h.StreamErr(e, err)
			return err
		}
		e.Summary = fmt.Sprintf("%s %s successfully", operation.Properties[config.OperationDescriptionKey], status)
		e.Details = fmt.Sprintf("%s is now %s.", operation.Properties[config.OperationDescriptionKey], status)
		h.StreamInfo(e)
	default:
		h.StreamErr(e, adapter.ErrOpInvalid)
	}
	return nil
}
