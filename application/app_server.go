package application

import (
	"gogen_grpc/domain_demogrpc/controller/graphqlserver"
	"gogen_grpc/domain_demogrpc/gateway/emptyimpl"
	"gogen_grpc/domain_demogrpc/usecase/runmessagereverse"
	"gogen_grpc/shared/config"
	"gogen_grpc/shared/gogen"
	"gogen_grpc/shared/infrastructure/logger"
)

type appServer struct{}

func NewAppServer() gogen.Runner {
	return &appServer{}
}

func (appServer) Run() error {

	const appName = "appserver"

	cfg := config.ReadConfig()

	appData := gogen.NewApplicationData(appName)

	log := logger.NewSimpleJSONLogger(appData)

	datasource := emptyimpl.NewGateway(log, appData, cfg)

	//primaryDriver := grpcserver.NewController(log, cfg)
	primaryDriver := graphqlserver.NewController(log, cfg)

	primaryDriver.AddUsecase(
		//
		runmessagereverse.NewUsecase(datasource),
	)

	primaryDriver.RegisterRouter()

	primaryDriver.Start()

	return nil
}
