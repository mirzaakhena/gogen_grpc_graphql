package application

import (
	"gogen_grpc/domain_demo/controller/restapi"
	"gogen_grpc/domain_demo/gateway/graphqlclient"
	"gogen_grpc/domain_demo/usecase/runmessagesend"
	"gogen_grpc/shared/config"
	"gogen_grpc/shared/gogen"
	"gogen_grpc/shared/infrastructure/logger"
	"gogen_grpc/shared/infrastructure/token"
)

type appClient struct{}

func NewAppClient() gogen.Runner {
	return &appClient{}
}

func (appClient) Run() error {

	const appName = "appClient"

	cfg := config.ReadConfig()

	appData := gogen.NewApplicationData(appName)

	log := logger.NewSimpleJSONLogger(appData)

	jwtToken := token.NewJWTToken(cfg.JWTSecretKey)

	//datasource := grpcclient.NewGateway(log, appData, cfg)
	datasource := graphqlclient.NewGateway(log, appData, cfg)

	primaryDriver := restapi.NewController(appData, log, cfg, jwtToken)

	primaryDriver.AddUsecase(
		//
		runmessagesend.NewUsecase(datasource),
	)

	primaryDriver.RegisterRouter()

	primaryDriver.Start()

	return nil
}
