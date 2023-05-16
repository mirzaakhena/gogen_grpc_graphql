package graphqlclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gogen_grpc/shared/config"
	"gogen_grpc/shared/gogen"
	"gogen_grpc/shared/infrastructure/logger"
	"io"
	"net/http"
)

type gateway struct {
	appData gogen.ApplicationData
	config  *config.Config
	log     logger.Logger
}

// NewGateway ...
func NewGateway(log logger.Logger, appData gogen.ApplicationData, cfg *config.Config) *gateway {

	return &gateway{
		log:     log,
		appData: appData,
		config:  cfg,
	}
}

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
