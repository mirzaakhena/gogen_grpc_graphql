package runmessagesend

import (
	"context"
)

type runMessageSendInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &runMessageSendInteractor{
		outport: outputPort,
	}
}

func (r *runMessageSendInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	returnMessage, err := r.outport.SendMessage(ctx, req.Message)
	if err != nil {
		return nil, err
	}

	res.ReturnMessage = returnMessage

	return res, nil
}
