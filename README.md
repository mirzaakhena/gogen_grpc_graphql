# Gogen grpc

In this repo we are demonstrating on how to use the grpc communication between application using the gogen framework

## Gogen Framework
For the Gogen Framework Structure, you can refer to here link

> https://github.com/mirzaakhena/gogen

## Application Architecture

The application consist of two parts
1. Client : Has a restapi interface to invoke the grpc client
2. Server : Has a grpc server to receive the request, process it and then return back to grpc client

![gogen grpc architecture](https://github.com/mirzaakhena/gogengrpc/blob/main/gogen_grpc_architecture.png)

## Folder structure
```text
gogen_grpc
├── application
│  ├── app_appclient.go
│  └── app_appserver.go
├── domain_demogrpc
│  ├── controller
│  │  ├── grpcreceiver
│  │  └── restapi
│  ├── gateway
│  │  ├── emptyimpl
│  │  └── grpcsender
│  └── usecase
│      ├── runmessagereverse
│      └── runmessagesend
├── main.go
└── shared
    └── pb
       ├── grpcstub
       │  ├── message.pb.go
       │  └── message_grpc.pb.go
       └── message.proto  
```

## How to run the application

1. After you git clone it, make sure to run the `go mod tidy` to download the dependency
2. Run the server application by `go run main.go appserver`
3. Run the client application by `go run main.go appclient`
4. invoke this api with curl, postman or use the file `http_runmessagesend.http` under `domain_demogrpc/controller/restapi`

    ```
    POST http://localhost:8000/api/v1/runmessagesend
    {
      "message": "hello" 
    }
    ```
    Then you will get the message reversed in response payload
    ```
   {
     "success": true,
     "errorCode": "",
     "errorMessage": "",
     "data": {
       "return_message": "olleh"
      },
     "traceId": "Z1RCGNXYTR2QCNVK"
   }   
    ```

## GRPC Stub Generation

This is the proto file `shared/pb/message.proto` used in this project 
```text
syntax = "proto3";

package mypackage;

option go_package = "./grpcstub";

message MessageReverseRequest {
  string content = 1;
}

message MessageReverseResponse {
  string content = 1;
}

service MyService {
  rpc SendMessage(MessageReverseRequest) returns (MessageReverseResponse);
}
```

For GRPC code generation you need to do
```
$ cd shared/pb
$ protoc --go_out=. --go-grpc_out=. message.proto
```
Make sure you already have the protoc executable first.

## GRPC Server in Controller

```go

type controller struct {
	grpcstub.UnimplementedMyServiceServer
	gogen.UsecaseRegisterer 
	server                  *grpc.Server
	log                     logger.Logger
	cfg                     *config.Config
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

func (r *controller) SendMessage(ctx context.Context, stubReq *grpcstub.MessageReverseRequest) (*grpcstub.MessageReverseResponse, error) {
	return &grpcstub.MessageReverseResponse{
		Content: "...",
	}, nil
}

```


## GRPC Client in Gateway

```go
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
	r.log.Info(ctx, "called")

	response, err := r.client.SendMessage(context.Background(), &grpcstub.MessageReverseRequest{
		Content: message,
	})
	if err != nil {
		return "", err
	}

	return response.Content, nil
}

```