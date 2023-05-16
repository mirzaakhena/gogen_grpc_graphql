# Gogen gRPC GraphQL

In this repo we are demonstrating on how to use the Grpc and GraphQL communication between application using the gogen framework

## Gogen Framework
For the Gogen Framework Structure, you can refer to here link

> https://github.com/mirzaakhena/gogen

## Application Architecture

The application consist of two parts
1. Client : Has a restapi interface to invoke the gRPC or GraphQL client
2. Server : Has a gRPC or GraphQL server to receive the request, process it and then return back to client (gRPC or GraphQL)

![gogen grpc architecture](https://github.com/mirzaakhena/gogengrpc/blob/main/gogen_grpc_architecture.png)

## Folder structure
```text
gogen_grpc
├── application
│  ├── app_client.go
│  └── app_server.go
├── domain_demo
│  ├── controller
│  │  ├── graphqlserver
│  │  ├── grpcserver
│  │  └── restapi
│  ├── gateway
│  │  ├── emptyimpl
│  │  ├── graphqlclient
│  │  └── grpcclient
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

## How to switch technology
For the client you may comment / uncomment this part
```
//primaryDriver := grpcserver.NewController(log, cfg)
primaryDriver := graphqlserver.NewController(log, cfg)
```

For the server you may comment / uncomment this part
```
//datasource := grpcclient.NewGateway(log, appData, cfg)
datasource := graphqlclient.NewGateway(log, appData, cfg)
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

## GRAPHQL Server in Controller

```
type controller struct {
	gogen.UsecaseRegisterer             // collect all the inports
	Router                  *gin.Engine // the router from preference web framework
	log                     logger.Logger
	cfg                     *config.Config
	fields                  graphql.Fields
}

func NewController(log logger.Logger, cfg *config.Config) gogen.ControllerRegisterer {

	return &controller{
		UsecaseRegisterer: gogen.NewBaseController(),
		log:               log,
		cfg:               cfg,
		fields:            map[string]*graphql.Field{},
	}

}

func (r *controller) RegisterRouter() {
	r.fields["reverseMessage"] = r.sendMessageHandler()
}

func (r *controller) Start() {

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: r.fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		fmt.Println("Error creating schema: ", err)
		return
	}

	// Create a new GraphQL HTTP handler with the schema
	graphqlHandler := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	// Serve the GraphQL endpoint
	http.Handle("/graphql", graphqlHandler)
	fmt.Println("GraphQL Server running on http://localhost:8080/graphql")
	http.ListenAndServe(":8080", nil)

}

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

```

## GRAPHQL Client in Gateway
```
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (r *gateway) SendMessage(ctx context.Context, message string) (string, error) {
	r.log.Info(ctx, "called in GraphQL Gateway")

	// Define the GraphQL query
	query := `
		query ReverseMessage($message: String!) {
			reverseMessage(message: $message)
		}
	`

	// Define the query variables
	variables := map[string]interface{}{
		"message": message,
	}

	// Create a GraphQL request
	request := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	// Convert the request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	// Send a POST request to the GraphQL server
	resp, err := http.Post("http://localhost:8080/graphql", "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the GraphQL response
	var response GraphQLResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	// Check for errors in the response
	if len(response.Errors) > 0 {
		errs := ""
		for _, err := range response.Errors {
			errs += err.Message + ", "
		}
		return "", fmt.Errorf(errs)
	}

	// Extract the reversed message from the response
	reversedMessage := response.Data.(map[string]interface{})["reverseMessage"].(string)

	return reversedMessage, nil
}
```