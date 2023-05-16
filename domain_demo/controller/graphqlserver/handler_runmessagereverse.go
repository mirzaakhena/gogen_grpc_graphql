package graphqlserver

import (
	"context"
	"fmt"
	"github.com/graphql-go/graphql"
	"gogen_grpc/domain_demo/usecase/runmessagereverse"
	"gogen_grpc/shared/gogen"
	"gogen_grpc/shared/infrastructure/logger"
	"gogen_grpc/shared/util"
)

func (r *controller) sendMessageHandler() *graphql.Field {

	return &graphql.Field{
		Type:        graphql.String,
		Description: "Reverses a given message",
		Args: graphql.FieldConfigArgument{
			"message": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {

			type InportRequest = runmessagereverse.InportRequest
			type InportResponse = runmessagereverse.InportResponse

			inport := gogen.GetInport[InportRequest, InportResponse](r.GetUsecase(InportRequest{}))

			traceID := util.GenerateID(16)

			ctx := logger.SetTraceID(context.Background(), traceID)

			var req InportRequest

			message, ok := p.Args["message"].(string)
			if !ok {
				return nil, fmt.Errorf("Invalid message type")
			}

			req.Message = message

			res, err := inport.Execute(ctx, req)
			if err != nil {
				return nil, err
			}

			return res.ReturnMessage, nil

		},
	}

}
