package grpcserver

import (
	"context"
	"gogen_grpc/domain_demo/usecase/runmessagereverse"
	"gogen_grpc/shared/gogen"
	"gogen_grpc/shared/pb/grpcstub"
)

func (r *controller) SendMessage(ctx context.Context, stubReq *grpcstub.MessageReverseRequest) (*grpcstub.MessageReverseResponse, error) {

	type InportRequest = runmessagereverse.InportRequest
	type InportResponse = runmessagereverse.InportResponse

	inport := gogen.GetInport[InportRequest, InportResponse](r.GetUsecase(InportRequest{}))

	var req InportRequest
	req.Message = stubReq.Content

	res, err := inport.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return &grpcstub.MessageReverseResponse{
		Content: res.ReturnMessage,
	}, nil

}
