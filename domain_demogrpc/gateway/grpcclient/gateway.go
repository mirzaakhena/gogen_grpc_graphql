package grpcclient

import (
	"context"
	"gogen_grpc/shared/config"
	"gogen_grpc/shared/gogen"
	"gogen_grpc/shared/infrastructure/logger"
	"gogen_grpc/shared/pb/grpcstub"
	"google.golang.org/grpc"
)

type gateway struct {
	appData gogen.ApplicationData
	config  *config.Config
	log     logger.Logger
	client  grpcstub.MyServiceClient
}

// NewGateway ...
func NewGateway(log logger.Logger, appData gogen.ApplicationData, cfg *config.Config) *gateway {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	//defer conn.Close()

	client := grpcstub.NewMyServiceClient(conn)

	return &gateway{
		log:     log,
		appData: appData,
		config:  cfg,
		client:  client,
	}
}

func (r *gateway) SendMessage(ctx context.Context, message string) (string, error) {
	r.log.Info(ctx, "called in GRPC Gateway")

	response, err := r.client.SendMessage(context.Background(), &grpcstub.MessageReverseRequest{
		Content: message,
	})
	if err != nil {
		return "", err
	}

	return response.Content, nil
}
