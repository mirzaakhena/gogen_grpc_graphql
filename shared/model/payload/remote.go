package payload

import "gogen_grpc/shared/gogen"

type Args struct {
	Type      string                `json:"type"`
	Data      any                   `json:"data"`
	Publisher gogen.ApplicationData `json:"publisher"`
	TraceID   string                `json:"traceId"`
}

type Reply struct {
	Success      bool
	ErrorMessage string
	Data         any
}
