package runmessagereverse

import (
	"gogen_grpc/shared/gogen"
)

type Inport = gogen.Inport[InportRequest, InportResponse]

type InportRequest struct {
	Message string
}

type InportResponse struct {
	ReturnMessage string
}
