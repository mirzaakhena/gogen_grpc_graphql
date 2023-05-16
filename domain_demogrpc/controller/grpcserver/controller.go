package grpcserver

import (
	"fmt"
	"gogen_grpc/shared/config"
	"gogen_grpc/shared/gogen"
	"gogen_grpc/shared/infrastructure/logger"
	"gogen_grpc/shared/pb/grpcstub"
	"google.golang.org/grpc"
	"net"
)

type controller struct {
	grpcstub.UnimplementedMyServiceServer
	gogen.UsecaseRegisterer
	server *grpc.Server
	log    logger.Logger
	cfg    *config.Config
}

func NewController(log logger.Logger, cfg *config.Config) gogen.ControllerRegisterer {

	server := grpc.NewServer()

	return &controller{
		UsecaseRegisterer: gogen.NewBaseController(),
		server:            server,
		log:               log,
		cfg:               cfg,
	}

}

func (r *controller) Start() {

	fmt.Println("GRPC Server is running on port 50051")

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	err = r.server.Serve(listen)
	if err != nil {
		panic(err)
	}
}

func (r *controller) RegisterRouter() {
	grpcstub.RegisterMyServiceServer(r.server, r)
}
