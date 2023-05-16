package runmessagereverse

import (
	"context"
)

type runMessageReverseInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &runMessageReverseInteractor{
		outport: outputPort,
	}
}

func (r *runMessageReverseInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	res.ReturnMessage = reverseString(req.Message)

	return res, nil
}

func reverseString(str string) string {
	runes := []rune(str)
	length := len(runes)
	for i := 0; i < length/2; i++ {
		runes[i], runes[length-1-i] = runes[length-1-i], runes[i]
	}
	return string(runes)
}
