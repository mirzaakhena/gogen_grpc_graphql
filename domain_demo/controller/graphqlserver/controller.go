package graphqlserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"gogen_grpc/shared/config"
	"gogen_grpc/shared/gogen"
	"gogen_grpc/shared/infrastructure/logger"
	"net/http"
)

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
