package payload

import (
	"gogen_grpc/shared/gogen"
)

type Payload struct {
	Data      any                   `json:"data"`
	Publisher gogen.ApplicationData `json:"publisher"`
	TraceID   string                `json:"traceId"`
}
